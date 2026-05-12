import type { RequestEvent } from "@sveltejs/kit";
import { json } from "@sveltejs/kit";

export async function GET({ url }: RequestEvent) {
	const id = url.searchParams.get('id');

	const data = {
		message: 'Hello from the API',
		requestedId: id
	};

	return json(data);
}
