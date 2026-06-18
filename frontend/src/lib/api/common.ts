import type { CfcCategory, User } from './model';

const COMMON_ENDPOINT = '/api/common';

export const userInfo = async (): Promise<User> => {
	const resp = await fetch('/api/auth');
	if (!resp.ok) {
		throw new Error('unknow error when asking for sign-in status');
	}

	const data = (await resp.json()) as { user: User };

	return data.user;
};

export const addCfc = async (
	subject: string,
	category: CfcCategory,
	details: string
): Promise<void> => {
	const resp = await fetch(`${COMMON_ENDPOINT}/cfc`, {
		method: 'POST',
		headers: new Headers({ 'Content-Type': 'application/json' }),
		body: JSON.stringify({ subject: subject, category: category, details: details })
	});

	if (!resp.ok) {
		throw new Error('error sending customer feedback/complaint');
	}
};
