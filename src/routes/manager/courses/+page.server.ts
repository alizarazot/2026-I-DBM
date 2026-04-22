import { auth } from "$lib/server/auth";
import { collectionCourses, type Course } from "$lib/server/database";
import { ObjectId } from "mongodb";
import type { PageData } from "../$types";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async (event) => {
  event.depends("manager:courses");
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
      id: course._id.toString(),
      name: course.name,
      description: course.description,
      day: course.day,
      startHour: course.startHour,
      duration: course.duration,
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
  getCourse: async (event) => {
    const formData = await event.request.formData();
    const course = await collectionCourses.findOne({
      _id: new ObjectId(formData.get("id")?.toString() ?? ""),
    });
    return { course: JSON.parse(JSON.stringify(course)) };
  },
  editCourse: async (event) => {
    const formData = await event.request.formData();
    const course: Course = {
      name: formData.get("name")?.toString() ?? "",
      description: formData.get("description")?.toString() ?? "",
      day: formData.get("day")?.toString() ?? "",
      startHour: formData.get("startHour")?.toString() ?? "",
      duration: formData.get("duration")?.toString() ?? "",
      maxStudents: parseInt(formData.get("maxStudents")?.toString() ?? "0"),
      teacherId: formData.get("teacherId")?.toString() ?? "",
    };
    const id = formData.get("id")?.toString() ?? "";
    console.log("id is", id);
    await collectionCourses.updateOne(
      { _id: new ObjectId(id) },
      { $set: course },
    );
  },
  deleteCourse: async (event) => {
    const formData = await event.request.formData();
    await collectionCourses.deleteOne({
      _id: new ObjectId(formData.get("id")?.toString() ?? ""),
    });
  },
  addCourse: async (event) => {
    const formData = await event.request.formData();
    const course: Course = {
      name: formData.get("name")?.toString() ?? "",
      description: formData.get("description")?.toString() ?? "",
      day: formData.get("day")?.toString() ?? "",
      startHour: formData.get("startHour")?.toString() ?? "",
      duration: formData.get("duration")?.toString() ?? "",
      maxStudents: parseInt(formData.get("maxStudents")?.toString() ?? "0"),
      teacherId: formData.get("teacherId")?.toString() ?? "",
    };
    await collectionCourses.insertOne(course);
  },
};
