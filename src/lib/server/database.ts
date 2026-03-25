import { env } from "$env/dynamic/private";

import { MongoClient } from "mongodb";

if (!env.MONGODB_URI) {
	throw new Error("Environment variable `MONGO_URI` is undefined!");
}
const client = new MongoClient(env.MONGODB_URI);

export const db = client.db("Plataforma-LINEA");

export type Course = {
	name: string;
	description: string;
	day: string;
	startHour: string;
	duration: number;
	maxStudents: number;
	teacherId: string;
};

export const collectionCourses = db.collection("courses");
