import { collectionAttendances } from '$lib/server/database';
import type { PageServerLoad } from '../$types';

export const load: PageServerLoad = async () => {
	const startOfToday = new Date();
	startOfToday.setHours(0, 0, 0, 0);

	const startOfTomorrow = new Date(startOfToday);
	startOfTomorrow.setDate(startOfTomorrow.getDate() + 1);

	const maleAttendances = await collectionAttendances
		.aggregate([
			{
				$match: {
					$and: [{ date: { $gte: startOfToday } }, { date: { $lt: startOfTomorrow } }]
				}
			},
			{
				$lookup: {
					from: 'users',
					localField: 'userId',
					foreignField: '_id',
					pipeline: [{ $project: { isMale: 1 } }],
					as: 'user_info'
				}
			},

			{ $unwind: '$user_info' },

			{
				$match: {
					$expr: {
						$eq: [{ $ifNull: ['$user_info.isMale', true] }, true]
					}
				}
			},

			{
				$group: {
					_id: {
						hour: { $hour: { date: '$date', timezone: 'America/Bogota' } },
						user: '$userId'
					}
				}
			},
			{
				$group: {
					_id: '$_id.hour',
					totalUniquePeople: { $sum: 1 }
				}
			},

			{
				$project: {
					_id: 0,
					formattedResult: {
						$arrayToObject: [[{ k: { $toString: '$_id' }, v: '$totalUniquePeople' }]]
					}
				}
			},
			{ $replaceRoot: { newRoot: '$formattedResult' } }
		])
		.next();

	const femaleAttendances = await collectionAttendances
		.aggregate([
			{
				$match: {
					$and: [{ date: { $gte: startOfToday } }, { date: { $lt: startOfTomorrow } }]
				}
			},
			{
				$lookup: {
					from: 'users',
					localField: 'userId',
					foreignField: '_id',
					pipeline: [{ $project: { isMale: 1 } }],
					as: 'user_info'
				}
			},

			{ $unwind: '$user_info' },

			{
				$match: {
					$expr: {
						$eq: ['$user_info.isMale', false]
					}
				}
			},

			{
				$group: {
					_id: {
						hour: { $hour: { date: '$date', timezone: 'America/Bogota' } },
						user: '$userId'
					}
				}
			},
			{
				$group: {
					_id: '$_id.hour',
					totalUniquePeople: { $sum: 1 }
				}
			},

			{
				$project: {
					_id: 0,
					formattedResult: {
						$arrayToObject: [[{ k: { $toString: '$_id' }, v: '$totalUniquePeople' }]]
					}
				}
			},
			{ $replaceRoot: { newRoot: '$formattedResult' } }
		])
		.next();

	return {
		maleAttendances: maleAttendances,
		femaleAttendances: femaleAttendances
	};
};
