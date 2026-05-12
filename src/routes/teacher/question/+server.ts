import { collectionQuestions } from '$lib/server/database';
import { json } from '@sveltejs/kit';

export const GET = async () => {
	const question = await collectionQuestions.find().sort({ _id: -1 }).limit(1).next();

	if (question?.lastAnswer == 0) {
		return json({ msg: 'Los estudiantes están prestando atención' });
	}

	if (!question?.lastAnswer) {
		return json({ msg: '' });
	}

	return json({ msg: question.badAnswers[question.lastAnswer - 1].teacherReinforcement });
};
