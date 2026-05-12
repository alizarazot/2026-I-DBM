import { json, type RequestHandler } from '@sveltejs/kit';
import { collectionUserAnalytics } from '$lib/server/database';

export const POST: RequestHandler = async ({ locals, request }) => {
	const body = await request.json();
	const { latitude, longitude } = body ?? {};

	if (typeof latitude !== 'number' || typeof longitude !== 'number') {
		return json({ error: 'Invalid location payload' }, { status: 400 });
	}

	await collectionUserAnalytics.insertOne({
		userId: locals.user.id,
		coordinates: [longitude, latitude],
		timestamp: new Date()
	});

	return json({ success: true });
};
