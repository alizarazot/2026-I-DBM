import { collectionCourses } from '$lib/server/database';
import type { Actions } from './$types';

const PAGING_SIZE = 3;

export const actions = {
	searchQuery: async ({ request }) => {
		const data = await request.formData();
		const searchQuery = (data.get('search-query') ?? '').toString();
		const isQueryInverted = (data.get('is-query-inverted') ?? 'false') == 'true' ? true : false;
		const pageQueryNumber = parseInt((data.get('page-query-number') ?? '0').toString());

		if (searchQuery.trim() == '') return { searchResults: [] };

		const queryResult = (
			!isQueryInverted
				? await collectionCourses
						.aggregate([
							{
								$match: {
									$text: { $search: searchQuery }
								}
							},
							{
								$skip: PAGING_SIZE * pageQueryNumber
							},
							{ $limit: PAGING_SIZE }
						])
						.toArray()
				: await collectionCourses
						.aggregate([
							{
								$match: {
									title: { $not: new RegExp(searchQuery, 'i') },
									description: { $not: new RegExp(searchQuery, 'i') }
								}
							},
							{
								$skip: PAGING_SIZE * pageQueryNumber
							},
							{ $limit: PAGING_SIZE }
						])
						.toArray()
		).map((result) => [result.name, result.description]);

		console.log(queryResult);

		return { searchResults: queryResult };
	}
} satisfies Actions;
