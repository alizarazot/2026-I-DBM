import { env } from "$env/dynamic/private";

import { MongoClient } from "mongodb";

if (!env.MONGODB_URI) {
	throw new Error("Environment variable `MONGO_URI` is undefined!");
}
const client = new MongoClient(env.MONGODB_URI);

export const db = client.db("Plataforma-LINEA");

export type Class = {
	transcription: [{ timestamp: number; text: string; personNumber: number }];
	date: Date;
	durationInSeconds: number;
};
export type Course = {
	name: string;
	description: string;
	courseKind: string;
	classIds: string[];
	maxStudents: number;
	teacherId: string;
	calendar: { day: string; hourStart: string; hourEnd: string }[];
};

export const collectionCourses = db.collection("courses");
