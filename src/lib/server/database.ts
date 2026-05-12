import { env } from '$env/dynamic/private';

import { MongoClient } from 'mongodb';

if (!env.MONGODB_URI) {
	throw new Error('Environment variable `MONGODB_URI` is undefined!');
}
const client = new MongoClient(env.MONGODB_URI);

export const db = client.db('Plataforma-LINEA');

export type UserLocationAnalytics = {
	userId: string;
	coordinates: number[]; // [longitude, latitude]
	timestamp: Date;
};

export type Lecture = {
	courseId: string;
	content: {
		userId: string;
		text: string;
		audio: string; // Base64-encoded audio data.
		timestamp: Date;
	}[];
	createdAt: Date;
};

export type Course = {
	name: string;
	description: string;
	maxStudents: number;
	studentsIds: string[];
	teacherId: string;
	createdAt: Date;
	schedules: {
		day: 'Sunday' | 'Monday' | 'Tuesday' | 'Wednesday' | 'Thursday' | 'Saturday';
		startTime: number; // Hour in 0000-2400 format.
		duration: number; // In minutes.
	}[];
};

export const collectionCourses = db.collection('courses');
export const collectionLectures = db.collection('lectures');
export const collectionQuestions = db.collection('questions');
export const collectionUserAnalytics = db.collection('user_location_analytics');
