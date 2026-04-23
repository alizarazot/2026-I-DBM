export async function GET({ url }) {
	// 1. Get query params if needed (e.g., ?id=123)
	const id = url.searchParams.get('id');

	const data = {
		message: 'Hello from the API',
		requestedId: id
	};

	// 2. Return using the json() helper
	return json(data);
}
