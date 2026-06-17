import type { User, UserRole } from '$lib/api/model';

const MANAGER_ENDPOINT = '/api/manager';

export const listUsers = async (
	page: number,
	limit: number,
	role: UserRole | null
): Promise<[User[], number]> => {
	const resp = await fetch(
		`${MANAGER_ENDPOINT}/users?${new URLSearchParams({ page: page.toString(), limit: limit.toString(), role: role ?? '' })}`
	);

	if (!resp.ok) {
		throw new Error('error while asking for user list');
	}

	const data = await resp.json();

	for (let i = 0; i < data.users.length; i++) {
		data.users[i].info.birthdate = new Date(data.users[i].info.birthdate);
	}

	return [data.users, data.totalUsers];
};

export const createUser = async (user: User, initialPassword: string): Promise<void> => {
	const resp = await fetch(`${MANAGER_ENDPOINT}/user`, {
		method: 'POST',
		headers: new Headers({ 'Content-Type': 'application/json' }),
		body: JSON.stringify({ user: user, initialPassword: initialPassword })
	});

	if (!resp.ok) {
		throw new Error('error registering user');
	}
};

export const updateUser = async (user: User): Promise<void> => {
	const resp = await fetch(`${MANAGER_ENDPOINT}/user`, {
		method: 'PUT',
		headers: new Headers({ 'Content-Type': 'application/json' }),
		body: JSON.stringify({ user: user })
	});

	if (!resp.ok) {
		throw new Error('error updating user');
	}
};

export const deleteUser = async (email: string): Promise<void> => {
	const resp = await fetch(`${MANAGER_ENDPOINT}/user?${new URLSearchParams({ email: email })}`, {
		method: 'DELETE'
	});

	if (!resp.ok) {
		throw new Error('error deleting user');
	}
};
