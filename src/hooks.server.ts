import type { Handle } from '@sveltejs/kit';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import { building } from '$app/environment';
import { auth } from '$lib/server/auth';

let initialized = false;

// TODO: Remove this, the user should be created on first seeding.
const handleBetterAuth: Handle = async ({ event, resolve }) => {
	if (!initialized) {
		initialized = true;
		try {
			await auth.api.createUser({
				body: {
					name: 'Manager',
					email: 'manager@local.app',
					password: 'manager@local.app',
					role: 'manager',
					data: {
						lastName: 'N/A',
						document: 'N/A',
						phone: 'N/A'
					}
				}
			});
			console.log('Root manager user created!');
		} catch (error) {
			console.log('Root manager user already exists, skipping creation');
		}
	}

	const session = await auth.api.getSession({ headers: event.request.headers });
	if (session) {
		event.locals.session = session.session;
		event.locals.user = session.user;
	} else {
		event.locals.session = null;
		event.locals.user = null;
	}

	return svelteKitHandler({ event, resolve, auth, building });
};

// This is the main entry point for handling requests in SvelteKit. We use better-auth's svelteKitHandler to manage authentication and session handling.
export const handle: Handle = handleBetterAuth;
