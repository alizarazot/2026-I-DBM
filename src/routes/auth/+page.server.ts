import { fail, redirect } from "@sveltejs/kit";
import { APIError } from "better-auth/api";
import { auth } from "$lib/server/auth";
import { db } from "$lib/server/database";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async (event: any) => {
	if (event.locals.user) {
		return redirect(302, "/");
	}
};

export const actions: Actions = {
	deleteUser: async (event) => {
		const formData = await event.request.formData();
		try {
			await auth.api.removeUser({
				body: {
					userId: formData.get("id")?.toString() ?? "",
				},
				headers: event.request.headers,
			});
		} catch (error) {
			if (error instanceof APIError) {
				return fail(400, { message: error.message || "Signin failed" });
			}
			return fail(500, { message: "Unexpected error" });
		}
	},
	updateUser: async (event) => {
		const formData = await event.request.formData();
		try {
			await auth.api.adminUpdateUser({
				body: {
					userId: formData.get("id")?.toString() ?? "",
					data: {
						name: formData.get("firstName")?.toString() ?? "",
						email: formData.get("email")?.toString() ?? "",
						lastName: formData.get("lastName")?.toString() ?? "",
						document: formData.get("document")?.toString() ?? "",
						phone: formData.get("phone")?.toString() ?? "",
					},
				},
				headers: event.request.headers,
			});
		} catch (error) {
			if (error instanceof APIError) {
				return fail(400, { message: error.message || "Signin failed" });
			}
			return fail(500, { message: "Unexpected error" });
		}
	},
	addManager: async (event) => {
		const formData = await event.request.formData();
		try {
			await auth.api.createUser({
				body: {
					name: formData.get("firstName")?.toString() ?? "",
					email: formData.get("email")?.toString() ?? "",
					password: formData.get("password")?.toString() ?? "",
					role: "manager",
					data: {
						lastName: formData.get("lastName")?.toString() ?? "",
						document: formData.get("document")?.toString() ?? "",
						phone: formData.get("phone")?.toString() ?? "",
					},
				},
			});
		} catch (error) {
			if (error instanceof APIError) {
				return fail(400, { message: error.message || "Signin failed" });
			}
			return fail(500, { message: "Unexpected error" });
		}
	},
	addTeacher: async (event) => {
		const formData = await event.request.formData();
		try {
			await auth.api.createUser({
				body: {
					name: formData.get("firstName")?.toString() ?? "",
					email: formData.get("email")?.toString() ?? "",
					password: formData.get("password")?.toString() ?? "",
					role: "teacher",
					data: {
						lastName: formData.get("lastName")?.toString() ?? "",
						document: formData.get("document")?.toString() ?? "",
						phone: formData.get("phone")?.toString() ?? "",
					},
				},
			});
		} catch (error) {
			if (error instanceof APIError) {
				return fail(400, { message: error.message || "Signin failed" });
			}
			return fail(500, { message: "Unexpected error" });
		}
	},
	addStudent: async (event) => {
		const formData = await event.request.formData();
		try {
			await auth.api.createUser({
				body: {
					name: formData.get("firstName")?.toString() ?? "",
					email: formData.get("email")?.toString() ?? "",
					password: formData.get("password")?.toString() ?? "",
					role: "student",
					data: {
						lastName: formData.get("lastName")?.toString() ?? "",
						document: formData.get("document")?.toString() ?? "",
						phone: formData.get("phone")?.toString() ?? "",
					},
				},
			});
		} catch (error) {
			if (error instanceof APIError) {
				return fail(400, { message: error.message || "Signin failed" });
			}
			return fail(500, { message: "Unexpected error" });
		}
	},
	signInEmail: async (event) => {
		const formData = await event.request.formData();
		const email = formData.get("email")?.toString() ?? "";
		const password = formData.get("password")?.toString() ?? "";

		try {
			await auth.api.signInEmail({
				body: {
					email,
					password,
				},
			});
		} catch (error) {
			if (error instanceof APIError) {
				return fail(400, { message: error.message || "Signin failed" });
			}
			return fail(500, { message: "Unexpected error" });
		}

		return redirect(302, "/");
	},
	signOut: async (event) => {
		await auth.api.signOut({
			headers: event.request.headers,
		});
		return redirect(302, "/auth");
	},
};
