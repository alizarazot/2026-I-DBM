import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { collectionCourses, collectionUserAnalytics, type Course } from '$lib/server/database';

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

	let currentDay: string;
	switch (new Date().getDay()) {
		case 0:
			currentDay = 'Sunday';
			break;
		case 1:
			currentDay = 'Monday';
			break;
		case 2:
			currentDay = 'Tuesday';
			break;
		case 3:
			currentDay = 'Wednesday';
			break;
		case 4:
			currentDay = 'Thursday';
			break;
		case 5:
			currentDay = 'Friday';
			break;
		case 6:
			currentDay = 'Saturday';
			break;
		default:
			throw new Error('Unexpected day of week');
	}
	const current1000Hour =
		new Date().getHours() * 100 + Math.floor(new Date().getMinutes() / 10) * 10;

	const courses = await collectionCourses
		.aggregate<Course>([
			{ $match: { studentsIds: locals.user.id } },
			{ $unwind: '$schedules' },
			{ $match: { 'schedules.day': currentDay } },
			{ $match: { 'schedules.startTime': { $lte: current1000Hour } } },
			{
				$match: {
					$expr: {
						$lt: [
							current1000Hour,
							{
								$add: ['$schedules.startTime', { $multiply: ['$schedules.duration', 100] }]
							}
						]
					}
				}
			},
			{ $limit: 1 }
		])
		.toArray();
	let currentCourse: Course | null = null;
	if (courses.length == 1) {
		currentCourse = courses[0];
	}

	return {
		user: locals.user,
		currentCourse: currentCourse,
		shouldAskLocation: shouldAskLocation
	};
};
