import { betterAuth } from "better-auth/minimal";
import { sveltekitCookies } from "better-auth/svelte-kit";
import { mongodbAdapter } from "better-auth/adapters/mongodb";
import { getRequestEvent } from "$app/server";
import { db } from "$lib/server/database";

export const auth = betterAuth({
	database: mongodbAdapter(db),
	emailAndPassword: { enabled: true },
	plugins: [sveltekitCookies(getRequestEvent)],
});
