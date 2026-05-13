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
		day: 'Sunday' | 'Monday' | 'Tuesday' | 'Wednesday' | 'Thursday' | 'Friday' | 'Saturday';
		startTime: number; // Hour in 0000-2400 format.
		duration: number; // In minutes.
	}[];
};

export type Attendance = {
	userId: string;
	date: Date;
};
export const collectionAttendances = db.collection('attendances');

export const collectionCourses = db.collection('courses');

export const weekdayNames = [
	'Sunday',
	'Monday',
	'Tuesday',
	'Wednesday',
	'Thursday',
	'Friday',
	'Saturday'
] as const;
export type WeekdayName = (typeof weekdayNames)[number];

export const getWeekdayName = (date = new Date()): WeekdayName => weekdayNames[date.getDay()];
export const getCurrent1000Hour = (date = new Date()) =>
	date.getHours() * 100 + Math.floor(date.getMinutes() / 10) * 10;

export const defaultCourse: Course = {
	name: 'General',
	description: 'General session for any unmatched course',
	maxStudents: 9999,
	studentsIds: [],
	teacherId: 'general',
	createdAt: new Date(0),
	schedules: []
};

const createCourseMatchPipeline = (matchFilter: Record<string, unknown>, date = new Date()) => {
	const currentDay = getWeekdayName(date);
	const current1000Hour = getCurrent1000Hour(date);

	return [
		{ $match: matchFilter },
		{ $unwind: '$schedules' },
		{ $match: { 'schedules.day': currentDay } },
		{ $match: { 'schedules.startTime': { $lte: current1000Hour } } },
		{
			$match: {
				$expr: {
					$lt: [
						current1000Hour,
						{
							$add: ['$schedules.startTime', { $multiply: ['$schedules.duration', 100] }]
						}
					]
				}
			}
		},
		{ $limit: 1 }
	];
};

const findCurrentCourse = async (filter: Record<string, unknown>, date = new Date()) => {
	const courses = await collectionCourses
		.aggregate<Course>(createCourseMatchPipeline(filter, date))
		.toArray();
	return courses.length === 1 ? courses[0] : defaultCourse;
};

export const getCurrentCourseForTeacher = async (teacherId: string, date = new Date()) =>
	findCurrentCourse({ teacherId }, date);

export const getCurrentCourseForStudent = async (studentId: string, date = new Date()) =>
	findCurrentCourse({ studentsIds: studentId }, date);

export type LatestUserLocation = {
	userId: string;
	coordinates: [number, number];
	timestamp: Date;
};

export const getLatestLocationsForUsers = async (userIds: string[]) => {
	if (!userIds || userIds.length === 0) {
		return [] as LatestUserLocation[];
	}

	return collectionUserAnalytics
		.aggregate<LatestUserLocation>([
			{ $match: { userId: { $in: userIds } } },
			{ $sort: { userId: 1, timestamp: -1 } },
			{
				$group: {
					_id: '$userId',
					userId: { $first: '$userId' },
					coordinates: { $first: '$coordinates' },
					timestamp: { $first: '$timestamp' }
				}
			}
		])
		.toArray();
};

export const getLatestLocationsForAllUsers = async () =>
	collectionUserAnalytics
		.aggregate<LatestUserLocation>([
			{ $sort: { userId: 1, timestamp: -1 } },
			{
				$group: {
					_id: '$userId',
					userId: { $first: '$userId' },
					coordinates: { $first: '$coordinates' },
					timestamp: { $first: '$timestamp' }
				}
			}
		])
		.toArray();

export const collectionLectures = db.collection('lectures');
export const collectionTranscriptions = db.collection('transcriptions');
export const collectionQuestions = db.collection('questions');
export const collectionUserAnalytics = db.collection('user_location_analytics');

export const getStudentCourseDetails = async (
	studentId: string,
	page: number = 1,
	limit: number = 10
) => {
	const skip = (page - 1) * limit;

	return collectionCourses
		.aggregate([
			{ $match: { studentsIds: studentId } },
			{
				$lookup: {
					from: 'users',
					localField: 'teacherId',
					foreignField: '_id',
					as: 'teacherDetails'
				}
			},
			{ $unwind: '$teacherDetails' },
			{
				$project: {
					_id: 0,
					courseName: '$name',
					teacherName: '$teacherDetails.firstName',
					teacherEmail: '$teacherDetails.email'
				}
			},
			{ $sort: { courseName: 1 } },
			{ $skip: skip },
			{ $limit: limit }
		])
		.toArray();
};

export const getLargeCoursesExceptGeneral = async (minStudents: number) => {
	return collectionCourses
		.find({
			maxStudents: { $gte: minStudents },
			name: { $ne: 'General' }
		})
		.toArray();
};

export const findUsersByCriteria = async (email: string) => {
	return db
		.collection('users')
		.find({
			$and: [
				{ email: { $not: { $eq: email } } },
				{
					$or: [{ role: 'teacher' }, { role: 'manager' }]
				}
			]
		})
		.toArray();
};
