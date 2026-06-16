<script lang="ts">
	import { Modal, Label, Input, Select, Datepicker, Button } from 'flowbite-svelte';

	import type { User, UserGenre, UserRole } from '$lib/api/model';

	const ROLES: { value: UserRole; name: string }[] = [
		{ value: 'manager', name: 'Manager' },
		{ value: 'teacher', name: 'Teacher' },
		{ value: 'student', name: 'Student' }
	];
	const GENRES: { value: UserGenre; name: string }[] = [
		{ value: 'male', name: 'Male' },
		{ value: 'female', name: 'Female' },
		{ value: 'other', name: 'Other' }
	];

	const INITIAL_FORM_DATA: FormData = {
		email: '',
		initialPassword: '',
		role: undefined,
		firstName: '',
		middleName: '',
		firstSurname: '',
		secondSurname: '',
		birthdate: undefined,
		genre: undefined
	};

	type FormData = {
		email: string;
		initialPassword: string;
		role: UserRole | undefined;
		firstName: string;
		middleName: string;
		firstSurname: string;
		secondSurname: string;
		birthdate: Date | undefined;
		genre: UserGenre | undefined;
	};

	let {
		user = $bindable(),
		showModal = $bindable(),
		onsave
	} = $props<{
		user: User | null;
		showModal: boolean;
		onsave: (user: User, initialPassword: string | null) => void;
	}>();

	const title = $derived(!user ? 'Add user' : 'Edit user');

	let formData = $state(INITIAL_FORM_DATA);

	const updateFormData = (): void => {
		if (!user) {
			formData = INITIAL_FORM_DATA;
			return;
		}

		formData = {
			email: user.email,
			initialPassword: '',
			role: user.role,
			firstName: user.info.firstName,
			middleName: user.info.middleName,
			firstSurname: user.info.firstSurname,
			secondSurname: user.info.secondSurname,
			birthdate: user.info.birthdate,
			genre: user.info.genre
		};
	};

	$effect(() => {
		if (showModal) {
			updateFormData();
		}
	});

	const handleSubmit = async (e: SubmitEvent): Promise<void> => {
		e.preventDefault();
		showModal = false;

		const newUser: User = {
			email: formData.email,
			role: formData.role ?? 'invalid',
			info: {
				firstName: formData.firstName,
				middleName: formData.middleName,
				firstSurname: formData.firstSurname,
				secondSurname: formData.secondSurname,
				birthdate:
					formData.birthdate ??
					(() => {
						throw new Error('invalid user birthdate');
					})(),
				genre: formData.genre ?? 'invalid'
			}
		};

		if (user) {
			onsave(newUser, null);
			return;
		}

		onsave(newUser, formData.initialPassword);
	};
</script>

<Modal {title} bind:open={showModal}>
	<form method="POST" action="/api/manager/add-user" onsubmit={handleSubmit}>
		<div class="mb-4 grid gap-4 sm:grid-cols-2">
			<div>
				<Label for="email" class="mb-2">Email</Label>
				<Input type="email" id="email" bind:value={formData.email} required />
			</div>
			<div>
				<Label>
					Role
					<Select id="genre" class="mt-2" items={ROLES} bind:value={formData.role} required />
				</Label>
			</div>
			<div>
				<Label for="first-name" class="mb-2">First name</Label>
				<Input type="text" id="first-name" bind:value={formData.firstName} required />
			</div>
			<div>
				<Label for="middle-name" class="mb-2">Middle name</Label>
				<Input type="text" id="middle-name" bind:value={formData.middleName} />
			</div>
			<div>
				<Label for="first-surname" class="mb-2">First surname</Label>
				<Input type="text" id="first-surname" bind:value={formData.firstSurname} required />
			</div>
			<div>
				<Label for="second-surname" class="mb-2">Second surname</Label>
				<Input type="text" id="second-surname" bind:value={formData.secondSurname} />
			</div>
			<div>
				<Label for="birthdate" class="mb-2">Birthdate</Label>
				<Datepicker id="birthdate" bind:value={formData.birthdate} required />
			</div>
			<div>
				<Label>
					Genre
					<Select id="genre" class="mt-2" items={GENRES} bind:value={formData.genre} required />
				</Label>
			</div>
			<div class="col-span-2">
				<Label for="initial-password" class="mb-2">Initial password</Label>
				<Input
					type="password"
					id="initial-password"
					bind:value={formData.initialPassword}
					required
				/>
			</div>
			<Button type="submit" class="w-52">
				<svg
					class="mr-1 -ml-1 h-6 w-6"
					fill="currentColor"
					viewBox="0 0 20 20"
					xmlns="http://www.w3.org/2000/svg"
					><path
						fill-rule="evenodd"
						d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z"
						clip-rule="evenodd"
					/></svg
				>
				{title}
			</Button>
		</div>
	</form>
</Modal>
