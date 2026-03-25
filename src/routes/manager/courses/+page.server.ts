import { collectionCourses, type Course } from "$lib/server/database";
import type { PageData } from "../$types";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async (event) => {
	const courses = [];
	for await (const course of collectionCourses.find()) {
		courses.push({
			name: course.name,
			description: course.description,
			courseKind: course.courseKind,
			maxStudents: course.maxStudents,
		});
	}

	return {
		courses: courses,
	};
};

export const actions: Actions = {
	addCourse: async (event) => {
		const formData = await event.request.formData();
		const course: Course = {
			name: formData.get("name")?.toString() ?? "",
			description: formData.get("description")?.toString() ?? "",
			courseKind: formData.get("courseKind")?.toString() ?? "",
			maxStudents: parseInt(formData.get("maxStudents")?.toString() ?? "0"),
			classIds: [],
			teacherId: "",
			calendar: [],
		};
		await collectionCourses.insertOne(course);
	},
};
