import { env } from '$env/dynamic/private';
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

	return json(await currentTranscription.json());
};
