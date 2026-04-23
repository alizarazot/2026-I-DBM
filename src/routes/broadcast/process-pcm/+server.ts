import { error, json } from "@sveltejs/kit";
import { auth } from "$lib/server/auth";
import type { RequestHandler } from "./$types";

let dataMap: Array<[number, string]> = [];

export const POST: RequestHandler = async ({ request }) => {
	const session = await auth.api.getSession({
		headers: request.headers,
	});

	if (!session) {
		throw error(401, "Unauthorized");
	}

	const data = await request.arrayBuffer();
	const transcription = await fetch("http://localhost:8000", {
		method: "POST",
		headers: {
			"Content-Type": "application/octet-stream",
			"X-Session-Id": session.session.id,
		},
		body: data,
	});

	const t = await transcription.json();

	if (!t.text || t.text.trim() === "") {
		return json({ text: t });
	}

	const lastItem = dataMap.length > 0 ? dataMap[dataMap.length - 1] : null;

	if (lastItem && lastItem[1] === t.text) {
		return json({ text: t });
	}

	dataMap.push([t.line, t.text]);

	if (dataMap.length > 20) {
		dataMap = dataMap.slice(-20);
	}

	return json({ text: t });
};

export const GET: RequestHandler = async ({ request }) => {
	return json(dataMap);
};
