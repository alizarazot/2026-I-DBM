import { collectionQuestions } from '$lib/server/database';
import { json } from '@sveltejs/kit';

export async function GET() {
	const questions = await collectionQuestions.find().sort({ _id: -1 }).limit(1).next();

	const lastQuestion = questions?.questions[questions.questions.length - 1];

	if (lastQuestion?.lastAnswer != null) {
		return json({ question: null });
	}

	return json(questions?.questions.pop());
}

export async function POST({ request }) {
	const data = await request.formData();
	const answer = parseInt(data.get('answer') as string);

	const questions = await collectionQuestions.find().sort({ _id: -1 }).limit(1).next();

	if (!questions) {
		return json({ success: false });
	}

	questions.questions[questions.questions.length - 1].lastAnswer = answer;

	collectionQuestions.updateOne(
		{ _id: questions!._id },
		{
			$set: {
				questions: questions!.questions,
				updatedAt: new Date()
			}
		}
	);

	return json({ success: true });
}
