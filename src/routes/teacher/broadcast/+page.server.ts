import type { PageServerLoad } from './$types';
import {
	getCurrentCourseForTeacher,
	getLatestLocationsForUsers,
	getLatestLocationsForAllUsers
} from '$lib/server/database';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.user) {
		throw redirect(302, '/login');
	}

	if (locals.user.role != 'teacher') {
		throw redirect(302, '/');
	}

	const currentCourse = await getCurrentCourseForTeacher(locals.user.id);
	const studentLocations =
		currentCourse.teacherId === 'general' || currentCourse.name === 'General'
			? await getLatestLocationsForAllUsers()
			: await getLatestLocationsForUsers(currentCourse.studentsIds);

	return {
		user: locals.user,
		currentCourse,
		studentLocations
	};
};
