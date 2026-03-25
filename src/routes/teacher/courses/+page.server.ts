import { collectionCourses } from "$lib/server/database";
import type { Actions } from "@sveltejs/kit";

export async function load({ locals }) {
	const courses = [];
	for await (const course of collectionCourses.find({
		teacherId: locals.user.id,
	})) {
		courses.push(JSON.parse(JSON.stringify(course)));
	}

	return {
		courses: courses,
	};
}
