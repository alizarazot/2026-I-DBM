import { error, json } from '@sveltejs/kit';
import { auth } from '$lib/server/auth';
import type { RequestHandler } from './$types';

const dataMap = [];

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

	const t = await transcription.json();

	try {
		let last = dataMap.pop();
		if (last[0] != t.line) {
			dataMap.push(last);
		}
	} catch (e) {}
	dataMap.push([t.line, t.text]);

	return json({ text: t });
};

export const GET: RequestHandler = async ({ request }) => {
	return json(dataMap);
};
