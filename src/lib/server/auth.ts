import { mongodbAdapter } from "better-auth/adapters/mongodb";
import { betterAuth } from "better-auth/minimal";
import { sveltekitCookies } from "better-auth/svelte-kit";
import { getRequestEvent } from "$app/server";
import { env } from "$env/dynamic/private";
import { db } from "$lib/server/database";

const baseURL = env.BETTER_AUTH_URL || "http://localhost:5173";

export const auth = betterAuth({
	baseURL,
	user: {
		modelName: "users",
		additionalFields: {
			role: {
				type: ["manager", "teacher", "student"],
				required: true,
			},
			document: {
				type: "string",
				required: true,
			},
			lastname: {
				type: "string",
				required: true,
			},
			phone: {
				type: "string",
				required: true,
			},
			isActive: {
				type: "boolean",
				required: true,
			},
			extraData: {
				type: "json",
				required: false,
			},
		},
	},
	session: {
		modelName: "user_sessions",
	},
	database: mongodbAdapter(db),
	emailAndPassword: { enabled: true },
	plugins: [sveltekitCookies(getRequestEvent)],
});
