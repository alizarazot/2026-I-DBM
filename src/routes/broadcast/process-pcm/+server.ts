import { env } from '$env/dynamic/private';
import { error, json } from '@sveltejs/kit';
import { auth } from '$lib/server/auth';

import type { RequestHandler } from './$types';

import { collectionQuestions, collectionTranscriptions } from '$lib/server/database';
import type { ObjectId } from 'mongodb';

import * as prompts from '$lib/prompts';

const TARGET_WORDS = 10;

export const POST: RequestHandler = async ({ request }) => {
	const session = await auth.api.getSession({
		headers: request.headers
	});

	if (!session) {
		throw error(401, 'Unauthorized');
	}

	const currentSessionId = request.headers.get('X-Current-Session-Id') ?? null;
	if (!currentSessionId) {
		throw error(401, 'Unauthorized');
	}

	const currentTranscription = await fetch(env.TTS_SERVER_URL ?? 'http://localhost:6061', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/octet-stream',
			'X-Session-Id': session.session.id,
			'X-Current-Session-Id': currentSessionId
		},
		body: await request.arrayBuffer()
	});

	const transcriptionLines = await currentTranscription.json();

	const result = await collectionTranscriptions.findOneAndUpdate(
		{ sessionId: `${session.session.id}-${currentSessionId}` },
		{
			$set: {
				lines: transcriptionLines,
				updatedAt: new Date()
			}
		},
		{ upsert: true, returnDocument: 'after' }
	);

	if (result) {
		await makeQuestions(result._id, transcriptionLines);
	}

	return json(transcriptionLines);
};

export const GET: RequestHandler = async ({}) => {
	const transcription = await collectionTranscriptions
		.find()
		.sort({ updatedAt: -1 })
		.limit(1)
		.next();

	if (!transcription) {
		return error(401, 'Not found');
	}

	return json(transcription.lines);
};

async function makeQuestions(id: ObjectId, lines: string[]) {
	let numWords = countWords(lines);

	const questionsData = await collectionQuestions.findOne(id);

	const lastPosition = questionsData?.lastPosition ?? 0;
	if (numWords < lastPosition + TARGET_WORDS) {
		return;
	}

	let oldQuestionsString = '';
	for (let question of questionsData?.questions ?? []) {
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
	for (let line of lines) {
		words += line.split(/\s+/).length;
	}
	return words;
}
