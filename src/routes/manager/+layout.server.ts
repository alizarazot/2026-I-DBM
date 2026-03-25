import { auth } from "$lib/server/auth";
import type { PageLoad } from "./$types";

export const load: PageLoad = async (event) => {
	try {
		const users = await auth.api.listUsers({
			query: {},
			headers: event.request.headers,
		});
		return {
			users: users.users,
		};
	} catch (err) {
		console.log(err);
	}
};
