<script lang="ts">
	import {
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		TableSearch,
		Button,
		Dropdown,
		ButtonGroup,
		List,
		Li,
		Radio,
		Input,
		Label,
		Modal,
		Select,
		Datepicker
	} from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import {
		PlusOutline,
		FilterSolid,
		ChevronRightOutline,
		ChevronLeftOutline
	} from 'flowbite-svelte-icons';
	import { onMount } from 'svelte';

	type UserRole = 'invalid' | 'manager' | 'teacher' | 'student';
	type UserGenre = 'invalid' | 'male' | 'female' | 'other';

	interface UserInfo {
		firstName: string;
		middleName: string;
		firstSurname: string;
		secondSurname: string;
		birthdate: Date;
		genre: UserGenre;
	}

	interface User {
		email: string;
		role: UserRole;
		info: UserInfo;
	}

	const usersPerPage = 10;

	let totalUsers = $state(0);
	let endPage = $derived(Math.floor(totalUsers / usersPerPage));

	let currentPage = $state(0);

	let pagesToShow = $derived.by(() => {
		const startPage = Math.min(1, currentPage);
		return Array.from({ length: endPage - startPage + 1 }, (_, i) => startPage + i);
	});

	let currentPageUsers: User[] = $state([]);

	let searchTerm = $state('');
	let filterRole = $state('');

	let filteredUsers = $derived(
		currentPageUsers.filter(
			(user) => user.email.toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1
		)
	);

	const updateDataAndPagination = async () => {
		const resp = await fetch(
			`/api/manager/list-users?${new URLSearchParams({ page: currentPage.toString(), limit: usersPerPage.toString(), role: filterRole })}`
		);
		const data = await resp.json();

		currentPageUsers = data.users;
		for (let i = 0; i < currentPageUsers.length; i++) {
			currentPageUsers[i].info.birthdate = new Date(currentPageUsers[i].info.birthdate);
		}

		totalUsers = data.totalUsers;
	};

	$effect(() => {
		// eslint-disable-next-line @typescript-eslint/no-unused-expressions
		filterRole; // To force update when [filterRole] changes.
		updateDataAndPagination();
	});

	const loadNextPage = () => {
		if (currentPage + usersPerPage < 0) {
			currentPage += usersPerPage;
			updateDataAndPagination();
		}
	};

	const loadPreviousPage = () => {
		if (currentPage - usersPerPage >= 0) {
			currentPage -= usersPerPage;
			updateDataAndPagination();
		}
	};

	const goToPage = (pageNumber: number) => {
		currentPage = (pageNumber - 1) * usersPerPage;
		updateDataAndPagination();
	};

	onMount(() => {
		updateDataAndPagination();
	});

	let showAddUserModal = $state(false);
	const handleSubmitAddUser = async (e: SubmitEvent) => {
		e.preventDefault();

		showAddUserModal = false;

		const target = e.target as HTMLFormElement;

		const userData: User = {
			email: addUserEmail,
			role: addUserRole ?? 'invalid',
			info: {
				firstName: addUserFirstName,
				middleName: addUserMiddleName,
				firstSurname: addUserFirstSurname,
				secondSurname: addUserSecondSurname,
				birthdate:
					addUserBirthdate ??
					(() => {
						throw new Error('invalid user birthdate');
					})(),
				genre: addUserGenre ?? 'invalid'
			}
		};

		const resp = await fetch(target.action, {
			method: target.method,
			headers: new Headers({ 'Content-Type': 'application/json' }),
			body: JSON.stringify(userData)
		});

		if (!resp.ok) {
			throw new Error('error registering user');
		}
	};

	const roles: { value: UserRole; name: string }[] = [
		{ value: 'manager', name: 'Manager' },
		{ value: 'teacher', name: 'Teacher' },
		{ value: 'student', name: 'Student' }
	];
	const genres: { value: UserGenre; name: string }[] = [
		{ value: 'male', name: 'Male' },
		{ value: 'female', name: 'Female' },
		{ value: 'other', name: 'Other' }
	];

	let addUserEmail = $state('');
	let addUserInitialPassword = $state('');
	let addUserRole: UserRole | undefined = $state(undefined);
	let addUserFirstName = $state('');
	let addUserMiddleName = $state('');
	let addUserFirstSurname = $state('');
	let addUserSecondSurname = $state('');
	let addUserBirthdate: Date | undefined = $state(undefined);
	let addUserGenre: UserGenre | undefined = $state(undefined);
</script>

<Section name="advancedTable" class="bg-gray-50 p-3 sm:p-5 dark:bg-gray-900">
	<TableSearch
		placeholder="Search by email"
		hoverable={true}
		classes={{
			root: 'bg-white dark:bg-gray-800 relative shadow-md sm:rounded-lg overflow-hidden',
			inner:
				'flex flex-col md:flex-row items-center justify-between space-y-3 md:space-y-0 md:space-x-4 p-4',
			search: 'w-full md:w-1/2 relative'
		}}
		tableClass="overflow-x-scroll block"
		bind:inputValue={searchTerm}
	>
		{#snippet header()}
			<div
				class="flex w-full flex-shrink-0 flex-col items-stretch justify-end space-y-2 md:w-auto md:flex-row md:items-center md:space-y-0 md:space-x-3"
			>
				<Button onclick={() => (showAddUserModal = true)}>
					<PlusOutline class="mr-2 h-3.5 w-3.5" />Add user
				</Button>
				<Button color="alternative">Filter<FilterSolid class="ml-2 h-3 w-3 " /></Button>
				<Dropdown class="w-48 space-y-2 p-3 text-sm">
					<h6 class="mb-3 text-sm font-medium text-gray-900 dark:text-white">Choose a role</h6>
					<List tag="dl">
						<Li>
							<Radio value="" bind:group={filterRole}>All roles</Radio>
						</Li>
						<Li>
							<Radio value="manager" bind:group={filterRole}>Managers</Radio>
						</Li>
						<Li>
							<Radio value="teacher" bind:group={filterRole}>Teachers</Radio>
						</Li>
						<Li>
							<Radio value="student" bind:group={filterRole}>Students</Radio>
						</Li>
					</List>
				</Dropdown>
			</div>
		{/snippet}
		<TableHead>
			<TableHeadCell class="px-4 py-3" scope="col">Email</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">First name</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">Middle name</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">First surname</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">Second surname</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">Role</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">Birthdate</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">Genre</TableHeadCell>
		</TableHead>
		<TableBody class="divide-y">
			{#if searchTerm !== ''}
				{#each filteredUsers as user (user.email)}
					<TableBodyRow>
						<TableBodyCell class="px-4 py-3">{user.email}</TableBodyCell>
						<TableBodyCell class="px-4 py-3">{user.info.firstName}</TableBodyCell>
						<TableBodyCell class="px-4 py-3">{user.info.middleName}</TableBodyCell>
						<TableBodyCell class="px-4 py-3">{user.info.firstSurname}</TableBodyCell>
						<TableBodyCell class="px-4 py-3">{user.info.secondSurname}</TableBodyCell>
						<TableBodyCell class="px-4 py-3 capitalize">{user.role}</TableBodyCell>
						<TableBodyCell class="px-4 py-3"
							>{user.info.birthdate.toLocaleDateString()}</TableBodyCell
						>
						<TableBodyCell class="px-4 py-3 capitalize">{user.info.genre}</TableBodyCell>
					</TableBodyRow>
				{/each}
			{:else}
				{#each currentPageUsers as user (user.email)}
					<TableBodyRow>
						<TableBodyCell class="px-4 py-3">{user.email}</TableBodyCell>
						<TableBodyCell class="px-4 py-3">{user.info.firstName}</TableBodyCell>
						<TableBodyCell class="px-4 py-3">{user.info.middleName}</TableBodyCell>
						<TableBodyCell class="px-4 py-3">{user.info.firstSurname}</TableBodyCell>
						<TableBodyCell class="px-4 py-3">{user.info.secondSurname}</TableBodyCell>
						<TableBodyCell class="px-4 py-3 capitalize">{user.role}</TableBodyCell>
						<TableBodyCell class="px-4 py-3"
							>{user.info.birthdate.toLocaleDateString()}</TableBodyCell
						>
						<TableBodyCell class="px-4 py-3 capitalize">{user.info.genre}</TableBodyCell>
					</TableBodyRow>
				{/each}
			{/if}
		</TableBody>
		{#snippet footer()}
			<div
				class="flex flex-col items-start justify-between space-y-3 p-4 md:flex-row md:items-center md:space-y-0"
				aria-label="Table navigation"
			>
				<span class="text-sm font-normal text-gray-500 dark:text-gray-400">
					{#if searchTerm == ''}
						Showing
					{:else}
						Searching on
					{/if}
					<span class="font-semibold text-gray-900 dark:text-white"
						>{currentPageUsers.length != 0
							? currentPage * usersPerPage + 1
							: 0}-{currentPageUsers.length}</span
					>
					of
					<span class="font-semibold text-gray-900 dark:text-white">{totalUsers}</span>
				</span>
				<ButtonGroup>
					<Button onclick={loadPreviousPage} disabled={currentPage === 0}
						><ChevronLeftOutline size="xs" class="m-1.5" /></Button
					>
					{#each pagesToShow as pageNumber (pageNumber)}
						<Button onclick={() => goToPage(pageNumber)}>{pageNumber}</Button>
					{/each}
					<Button onclick={loadNextPage} disabled={currentPage === endPage}
						><ChevronRightOutline size="xs" class="m-1.5" /></Button
					>
				</ButtonGroup>
			</div>
		{/snippet}
	</TableSearch>
</Section>

<Modal title="Add User" bind:open={showAddUserModal}>
	<form method="POST" action="/api/manager/add-user" onsubmit={handleSubmitAddUser}>
		<div class="mb-4 grid gap-4 sm:grid-cols-2">
			<div>
				<Label for="email" class="mb-2">Email</Label>
				<Input type="email" id="email" bind:value={addUserEmail} required />
			</div>
			<div>
				<Label>
					Role
					<Select id="genre" class="mt-2" items={roles} bind:value={addUserRole} required />
				</Label>
			</div>
			<div>
				<Label for="first-name" class="mb-2">First name</Label>
				<Input type="text" id="first-name" bind:value={addUserFirstName} required />
			</div>
			<div>
				<Label for="middle-name" class="mb-2">Middle name</Label>
				<Input type="text" id="middle-name" bind:value={addUserMiddleName} />
			</div>
			<div>
				<Label for="first-surname" class="mb-2">First surname</Label>
				<Input type="text" id="first-surname" bind:value={addUserFirstSurname} required />
			</div>
			<div>
				<Label for="second-surname" class="mb-2">Second surname</Label>
				<Input type="text" id="second-surname" bind:value={addUserSecondSurname} />
			</div>
			<div>
				<Label for="birthdate" class="mb-2">Birthdate</Label>
				<Datepicker id="birthdate" bind:value={addUserBirthdate} required />
			</div>
			<div>
				<Label>
					Genre
					<Select id="genre" class="mt-2" items={genres} bind:value={addUserGenre} required />
				</Label>
			</div>
			<div class="col-span-2">
				<Label for="initial-password" class="mb-2">Initial password</Label>
				<Input type="password" id="initial-password" bind:value={addUserInitialPassword} required />
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
				Add user
			</Button>
		</div>
	</form>
</Modal>
