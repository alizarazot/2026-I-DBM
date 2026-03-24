import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { auth } from "$lib/server/auth";

export const load: PageServerLoad = async (event) => {
	if (event.url.pathname === "/auth") {
		return;
	}
	if (!event.locals.user) {
		return redirect(302, "/auth");
	}
	return { user: event.locals.user };
};
