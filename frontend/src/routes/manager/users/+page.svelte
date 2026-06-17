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
		Radio
	} from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import {
		PlusOutline,
		FilterSolid,
		ChevronRightOutline,
		ChevronLeftOutline
	} from 'flowbite-svelte-icons';

	import UserModalForm from '$lib/components/manager/UserModalForm.svelte';

	import { createUser, listUsers, updateUser, deleteUser } from '$lib/api/manager';
	import type { User, UserRole } from '$lib/api/model';

	const USERS_PER_PAGE = 8;
	const PAGINATION_PADDING = 3;

	let filterRole: UserRole | null = $state(null);
	let currentPageNumber = $state(0);

	let currentPageUsers: User[] = $state([]);
	let totalUsers = $state(0);

	const updateData = async (
		pageNumber: number,
		filterRole: UserRole | null
	): Promise<[User[], number]> => {
		return await listUsers(pageNumber, USERS_PER_PAGE, filterRole);
	};

	$effect(() => {
		updateData(currentPageNumber, filterRole).then(([users, total]) => {
			currentPageUsers = users;
			totalUsers = total;
		});
	});

	let searchTerm = $state('');

	let lastPage = $derived(Math.floor(totalUsers / USERS_PER_PAGE) - 1);
	let pageNumbersToShow = $derived.by(() => {
		const startPage = Math.max(1, currentPageNumber - PAGINATION_PADDING);
		const endPage = Math.min(lastPage + 1, currentPageNumber + PAGINATION_PADDING);
		return Array.from({ length: endPage - startPage + 1 }, (_, i) => startPage + i);
	});

	let filteredUsers = $derived(
		currentPageUsers.filter(
			(user) => user.email.toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1
		)
	);

	const renderUsers = $derived(searchTerm != '' ? filteredUsers : currentPageUsers);

	const loadNextPage = () => {
		currentPageNumber++;
	};
	const loadPreviousPage = () => {
		currentPageNumber--;
	};
	const goToPage = (pageNumber: number) => {
		currentPageNumber = pageNumber - 1;
	};

	let modalUser = $state<User | null>(null);
	let showModal = $state(false);

	const addUser = (): void => {
		modalUser = null;
		showModal = true;
	};

	const editUser = (user: User): void => {
		modalUser = user;
		showModal = true;
	};

	const handlerDeleteUser = async (email: string): Promise<void> => {
		await deleteUser(email);
		[currentPageUsers, totalUsers] = await updateData(currentPageNumber, filterRole);
	};

	const handleModalOnSave = async (user: User, initialPassword: string | null): Promise<void> => {
		if (initialPassword) {
			await createUser(user, initialPassword);
		} else {
			await updateUser(user);
		}

		[currentPageUsers, totalUsers] = await updateData(currentPageNumber, filterRole);
	};
</script>

<Section name="advancedTable" class="bg-gray-50 p-3 sm:p-5 dark:bg-gray-900">
	<TableSearch
		placeholder="Search by email"
		hoverable={true}
		classes={{
			root: 'bg-white dark:bg-gray-800 relative shadow-md sm:rounded-lg overflow-hidden',
			inner:
				'flex flex-col md:flex-row items-center justify-between space-y-3 md:space-y-0 md:space-x-4 p-4',
			search: 'w-full md:w-1/2 relative',
			table: 'overflow-x-scroll block'
		}}
		bind:inputValue={searchTerm}
	>
		{#snippet header()}
			<div
				class="flex w-full flex-shrink-0 flex-col items-stretch justify-end space-y-2 md:w-auto md:flex-row md:items-center md:space-y-0 md:space-x-3"
			>
				<Button onclick={addUser}>
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
			<TableHeadCell />
			<TableHeadCell />
		</TableHead>
		<TableBody class="divide-y">
			{#each renderUsers as user (user.email)}
				<TableBodyRow>
					<TableBodyCell class="px-4 py-3">{user.email}</TableBodyCell>
					<TableBodyCell class="px-4 py-3">{user.info.firstName}</TableBodyCell>
					<TableBodyCell class="px-4 py-3">{user.info.middleName}</TableBodyCell>
					<TableBodyCell class="px-4 py-3">{user.info.firstSurname}</TableBodyCell>
					<TableBodyCell class="px-4 py-3">{user.info.secondSurname}</TableBodyCell>
					<TableBodyCell class="px-4 py-3 capitalize">{user.role}</TableBodyCell>
					<TableBodyCell class="px-4 py-3">{user.info.birthdate.toLocaleDateString()}</TableBodyCell
					>
					<TableBodyCell class="px-4 py-3 capitalize">{user.info.genre}</TableBodyCell>
					<TableBodyCell class="px-4 py-3"
						><Button
							onclick={() => {
								editUser(user);
							}}>Edit</Button
						></TableBodyCell
					>
					<TableBodyCell class="px-4 py-3"
						><Button
							color="red"
							onclick={() => {
								handlerDeleteUser(user.email);
							}}>Delete</Button
						></TableBodyCell
					>
				</TableBodyRow>
			{/each}
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
							? currentPageNumber * USERS_PER_PAGE + 1
							: 0}-{currentPageUsers.length}</span
					>
					of
					<span class="font-semibold text-gray-900 dark:text-white">{totalUsers}</span>
				</span>
				<ButtonGroup>
					<Button onclick={loadPreviousPage} disabled={currentPageNumber === 0}
						><ChevronLeftOutline size="xs" class="m-1.5" /></Button
					>
					{#each pageNumbersToShow as pageNumber (pageNumber)}
						<Button onclick={() => goToPage(pageNumber)}>{pageNumber}</Button>
					{/each}
					<Button onclick={loadNextPage} disabled={currentPageNumber === lastPage}
						><ChevronRightOutline size="xs" class="m-1.5" /></Button
					>
				</ButtonGroup>
			</div>
		{/snippet}
	</TableSearch>
</Section>

<UserModalForm bind:showModal onsave={handleModalOnSave} bind:user={modalUser} />
