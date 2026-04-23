import { auth } from '$lib/server/auth';
import { collectionCourses } from '$lib/server/database';
import type { Actions } from '@sveltejs/kit';
import { ObjectId } from 'mongodb';
import type { PageServerLoad } from '../auth/$types';

export const load: PageServerLoad = async (event) => {
	event.depends('manager:users');
	try {
		const users = await auth.api.listUsers({
			query: {},
			headers: event.request.headers
		});
		return {
			users: users.users
		};
	} catch (err) {
		console.log(err);
	}
};
