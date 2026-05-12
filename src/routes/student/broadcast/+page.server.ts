import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { collectionUserAnalytics, getCurrentCourseForStudent } from '$lib/server/database';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.user) {
		throw redirect(302, '/login');
	}

	if (locals.user.role != 'student') {
		throw redirect(302, '/');
	}

	// Check if user location should be asked again (every 5 minutes).
	const lastLocation = await collectionUserAnalytics
		.aggregate([
			{ $match: { userId: locals.user.id } },
			{ $unwind: '$locations' },
			{ $sort: { 'locations.timestamp': -1 } },
			{ $limit: 1 }
		])
		.toArray();

	let shouldAskLocation = true;
	if (lastLocation.length == 1) {
		const lastTimestamp = lastLocation[0].locations.timestamp;
		const now = new Date();
		const diffMinutes = (now.getTime() - lastTimestamp.getTime()) / 60000;
		shouldAskLocation = diffMinutes >= 5;
	}

	const currentCourse = await getCurrentCourseForStudent(locals.user.id);

	return {
		user: locals.user,
		currentCourse,
		shouldAskLocation
	};
};
