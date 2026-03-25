import type { Handle } from "@sveltejs/kit";
import { svelteKitHandler } from "better-auth/svelte-kit";
import { building } from "$app/environment";
import { auth } from "$lib/server/auth";

let initialized = false;

const handleBetterAuth: Handle = async ({ event, resolve }) => {
	if (!initialized) {
		initialized = true;
		try {
			await auth.api.createUser({
				body: {
					name: "Manager (Localhost)",
					email: "manager@localhost.app",
					password: "manager@localhost.app",
					role: "admin",
					data: {
						lastname: "",
						subrole: "manager",
						document: "",
						phone: "",
						isActive: true,
					},
				},
			});
			console.log("Root manager user created!");
		} catch (error) {
			console.log("Root manager user already exists, skipping creation");
		}
	}

	const session = await auth.api.getSession({ headers: event.request.headers });

	if (session) {
		event.locals.session = session.session;
		event.locals.user = session.user;
	}

	return svelteKitHandler({ event, resolve, auth, building });
};

export const handle: Handle = handleBetterAuth;
