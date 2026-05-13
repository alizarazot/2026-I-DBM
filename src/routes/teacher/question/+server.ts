import { collectionQuestions } from '$lib/server/database';
import { json } from '@sveltejs/kit';

export const GET = async () => {
	const question = await collectionQuestions.find().sort({ _id: -1 }).limit(1).next();

	if (!question || question.updatedAt?.getTime() < new Date().getTime() - 300000) {
		return json({ msg: '' });
	}

	const lastQuestion = question.questions[question.questions.length - 1];

	if (lastQuestion?.lastAnswer < 0) {
		return json({ msg: 'Los estudiantes no están prestando atención' });
	}

	if (lastQuestion?.lastAnswer == 0) {
		return json({ msg: 'El tema se está explicando de forma correcta' });
	}

	if (!lastQuestion?.lastAnswer) {
		return json({ msg: '' });
	}

	return json({ msg: lastQuestion.badAnswers[lastQuestion.lastAnswer - 1].teacherReinforcement });
};
