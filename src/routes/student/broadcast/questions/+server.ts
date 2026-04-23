import { collectionQuestions } from '$lib/server/database';
import { json } from '@sveltejs/kit';

export async function GET({}) {
	const questions = await collectionQuestions.find().sort({ _id: -1 }).limit(1).next();

	return json(questions?.questions.pop());
}
