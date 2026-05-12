import type { PageServerLoad } from './$types';
import { collectionCourses, type Course } from '$lib/server/database';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.user) {
		throw redirect(302, '/login');
	}

	if (locals.user.role != 'teacher') {
		throw redirect(302, '/');
	}

	const now = new Date();

	let currentDay: string;
	switch (now.getDay()) {
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

	const current1000Hour = now.getHours() * 100 + Math.floor(now.getMinutes() / 10) * 10;

	const courses = await collectionCourses
		.aggregate<Course>([
			{ $match: { teacherId: locals.user.id } },
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
		currentCourse: currentCourse
	};
};
