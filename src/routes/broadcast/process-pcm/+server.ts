import { error, json } from '@sveltejs/kit';
import { auth } from '$lib/server/auth';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request }) => {
	const session = await auth.api.getSession({
		headers: request.headers
	});

	if (!session) {
		throw error(401, 'Unauthorized');
	}

	const data = await request.arrayBuffer();
	const transcription = await fetch('http://localhost:8000', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/octet-stream',
			'X-Session-Id': session.session.id
		},
		body: data
	});

	return json({ text: await transcription.json() });
};
