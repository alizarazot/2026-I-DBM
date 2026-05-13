import { env } from '$env/dynamic/private';
import { error, json } from '@sveltejs/kit';
import { auth } from '$lib/server/auth';

import type { RequestHandler } from './$types';

import { collectionQuestions, collectionTranscriptions } from '$lib/server/database';
import type { ObjectId } from 'mongodb';
import { SpeechClient } from '@google-cloud/speech/build/src';

import * as prompts from '$lib/prompts';

const TARGET_WORDS = 10;
const SAMPLE_RATE = 16000;
const LANGUAGE_CODE = env.GOOGLE_SPEECH_LANGUAGE_CODE ?? 'es-ES';
const speechClient = new SpeechClient();

const float32To16BitPCM = (float32Array: Float32Array) => {
	const buffer = new ArrayBuffer(float32Array.length * 2);
	const view = new DataView(buffer);

	for (let i = 0; i < float32Array.length; i++) {
		let sample = Math.max(-1, Math.min(1, float32Array[i]));
		sample = sample < 0 ? sample * 0x8000 : sample * 0x7fff;
		view.setInt16(i * 2, sample, true);
	}

	return Buffer.from(buffer);
};

export const POST: RequestHandler = async ({ request }) => {
	const session = await auth.api.getSession({
		headers: request.headers
	});

	if (!session) {
		throw error(401, 'Unauthorized');
	}

	const rawAudio = await request.arrayBuffer();
	const float32Audio = new Float32Array(rawAudio);
	const pcmAudio = float32To16BitPCM(float32Audio);

	const audioBytes = pcmAudio.toString('base64');

	const [response] = await speechClient.recognize({
		audio: {
			content: audioBytes
		},
		config: {
			encoding: 'LINEAR16',
			sampleRateHertz: SAMPLE_RATE,
			languageCode: LANGUAGE_CODE,
			enableAutomaticPunctuation: true
		}
	});

	const transcriptionLines =
		response.results?.flatMap((result) =>
			result?.alternatives?.[0]?.transcript ? [result.alternatives[0].transcript] : []
		) ?? [];

	const updateOps: Record<string, unknown> =
		transcriptionLines.length > 0
			? {
					$set: {
						updatedAt: new Date()
					},
					$push: {
						lines: {
							$each: transcriptionLines
						}
					}
				}
			: {
					$set: {
						updatedAt: new Date()
					},
					$setOnInsert: {
						lines: []
					}
				};

	await collectionTranscriptions.findOneAndUpdate({ sessionId: session.session.id }, updateOps, {
		upsert: true
	});

	const result = await collectionTranscriptions.findOne({ sessionId: session.session.id });

	const lines = result?.lines ?? [];
	if (lines.length > 0) {
		await makeQuestions(result!._id, lines);
	}

	return json({ lines });
};

export const GET: RequestHandler = async () => {
	const transcription = await collectionTranscriptions
		.find()
		.sort({ updatedAt: -1 })
		.limit(1)
		.next();
	if (!transcription) {
		throw error(404, 'Not found');
	}

	return json({ lines: transcription.lines ?? [] });
};

async function makeQuestions(id: ObjectId, lines: string[]) {
	const numWords = countWords(lines);

	const questionsData = await collectionQuestions.findOne(id);

	const lastPosition = questionsData?.lastPosition ?? 0;
	if (numWords < lastPosition + TARGET_WORDS) {
		return;
	}

	let oldQuestionsString = '';
	for (const question of questionsData?.questions ?? []) {
		oldQuestionsString += ` - ${question.question}\n`;
	}
	const body = JSON.stringify({
		model: 'gemma3:12b-cloud', // TODO: Take from environment variable.
		system: prompts.SYSTEM_PROMPT_QUESTION_MAKER,
		prompt: prompts.format(prompts.CONTENT_PROMPT_QUESTION_MAKER, {
			content: lines.join(' '),
			questions: oldQuestionsString
		}),
		stream: false,
		format: 'json'
	});

	// TODO: Use env var for endpoint.
	const req = await fetch('http://localhost:11434/api/generate', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: body
	});

	try {
		// Take only the JSON from the LLM answer.
		const question = JSON.parse(((await req.json()).response.match(/\{[\s\S]*\}/) ?? [''])[0]);

		question.lastAnswer = null;

		collectionQuestions.updateOne(
			{ _id: id },
			{
				$push: {
					questions: question
				},
				$set: {
					lastPosition: numWords
				}
			},
			{ upsert: true }
		);
	} catch (e) {
		console.error(e);
	}
}

function countWords(lines: string[]): number {
	let words = 0;
	for (const line of lines) {
		words += line.split(/\s+/).length;
	}
	return words;
}
