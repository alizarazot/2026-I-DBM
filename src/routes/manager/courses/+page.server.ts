import { auth } from "$lib/server/auth";
import { collectionCourses, type Course } from "$lib/server/database";
import type { PageData } from "../$types";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async (event) => {
	let users: any;
	try {
		users = await auth.api.listUsers({
			query: {},
			headers: event.request.headers,
		});
	} catch (err) {
		console.log(err);
	}

	const courses = [];
	for await (const course of collectionCourses.find()) {
		courses.push({
			name: course.name,
			description: course.description,
			maxStudents: course.maxStudents,
			teacherId: course.teacherId,
		});
	}

	return {
		courses: courses,
		users: users.users,
	};
};

export const actions: Actions = {
	addCourse: async (event) => {
		const formData = await event.request.formData();
		const course: Course = {
			name: formData.get("name")?.toString() ?? "",
			description: formData.get("description")?.toString() ?? "",
			maxStudents: parseInt(formData.get("maxStudents")?.toString() ?? "0"),
			teacherId: formData.get("teacherId")?.toString() ?? "",
		};
		await collectionCourses.insertOne(course);
	},
};
