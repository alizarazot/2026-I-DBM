import { mongodbAdapter } from "better-auth/adapters/mongodb";
import { betterAuth } from "better-auth/minimal";
import { admin } from "better-auth/plugins";
import { createAccessControl } from "better-auth/plugins/access";
import { defaultStatements, adminAc } from "better-auth/plugins/admin/access";
import { sveltekitCookies } from "better-auth/svelte-kit";
import { getRequestEvent } from "$app/server";
import { env } from "$env/dynamic/private";
import { db } from "$lib/server/database";

const baseURL = env.BETTER_AUTH_URL || "http://localhost:5173";

const statement = {
	...defaultStatements,
	project: [],
} as const;
const ac = createAccessControl(statement);

const manager = ac.newRole({ ...adminAc.statements });
const teacher = ac.newRole({ project: [] });
const student = ac.newRole({ project: [] });

export const auth = betterAuth({
	baseURL,
	user: {
		modelName: "users",
		fields: {
			name: "fullName",
		},
		additionalFields: {
			document: {
				type: "string",
				required: true,
			},
			lastName: {
				type: "string",
				required: true,
			},
			phone: {
				type: "string",
				required: true,
			},
		},
	},
	session: {
		modelName: "user_sessions",
	},
	account: {
		modelName: "user_accounts",
	},
	database: mongodbAdapter(db),
	emailAndPassword: { enabled: true },
	plugins: [
		admin({
			ac,
			adminRoles: ["manager"],
			roles: {
				manager,
				teacher,
				student,
			},
		}),
		sveltekitCookies(getRequestEvent),
	],
});
