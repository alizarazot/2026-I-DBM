import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./auth/$types";

export const load: PageServerLoad = async (event) => {
	if (event.url.pathname === "/auth") {
		return;
	}
	if (!event.locals.user) {
		return redirect(302, "/auth");
	}

	if (event.url.pathname === "/") {
		switch (event.locals.user.subrole) {
			case "manager":
				return redirect(302, "/manager");
			case "teacher":
				return redirect(302, "/teacher");
			case "student":
				return redirect(302, "/student");
		}
	}

	return { user: event.locals.user };
};
