import { auth } from '$lib/server/auth';
import type { Actions } from '@sveltejs/kit';

export const actions: Actions = {
	getUser: async (event) => {
		const formData = await event.request.formData();
		const id = formData.get('id')?.toString() ?? '';

		const data = await auth.api.getUser({
			query: {
				id: id
			},
			headers: event.request.headers
		});

		return {
			user: data
		};
	}
};
