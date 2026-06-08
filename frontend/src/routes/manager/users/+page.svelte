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
		DropdownItem,
		Checkbox,
		ButtonGroup,
		List,
		Li
	} from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import {
		PlusOutline,
		ChevronDownOutline,
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

	let filteredUsers = $derived(
		currentPageUsers.filter(
			(user) => user.email.toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1
		)
	);

	const updateDataAndPagination = async () => {
		const respPage = await fetch(
			`/api/manager/users?${new URLSearchParams({ page: currentPage.toString(), limit: usersPerPage.toString() })}`
		);
		currentPageUsers = (await respPage.json()).users;

		const respTotal = await fetch('/api/manager/users/total');
		totalUsers = (await respTotal.json()).total;
	};

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
</script>

<Section name="advancedTable" sectionClass="bg-gray-50 dark:bg-gray-900 p-3 sm:p-5">
	<TableSearch
		placeholder="Search"
		hoverable={true}
		classes={{
			root: 'bg-white dark:bg-gray-800 relative shadow-md sm:rounded-lg overflow-hidden',
			inner:
				'flex flex-col md:flex-row items-center justify-between space-y-3 md:space-y-0 md:space-x-4 p-4',
			search: 'w-full md:w-1/2 relative'
		}}
		bind:inputValue={searchTerm}
	>
		{#snippet header()}
			<div
				class="flex w-full flex-shrink-0 flex-col items-stretch justify-end space-y-2 md:w-auto md:flex-row md:items-center md:space-y-0 md:space-x-3"
			>
				<Button>
					<PlusOutline class="mr-2 h-3.5 w-3.5" />Add user
				</Button>
				<Button color="alternative">Actions<ChevronDownOutline class="ml-2 h-3 w-3 " /></Button>
				<Dropdown simple class="w-44 divide-y divide-gray-100">
					<DropdownItem>Mass Edit</DropdownItem>
					<DropdownItem>Delete all</DropdownItem>
				</Dropdown>
				<Button color="alternative">Filter<FilterSolid class="ml-2 h-3 w-3 " /></Button>
				<Dropdown class="w-48 space-y-2 p-3 text-sm">
					<h6 class="mb-3 text-sm font-medium text-gray-900 dark:text-white">Choose brand</h6>
					<List tag="dl">
						<Li>
							<Checkbox>Apple (56)</Checkbox>
						</Li>
						<Li>
							<Checkbox>Microsoft (16)</Checkbox>
						</Li>
						<Li>
							<Checkbox>Razor (49)</Checkbox>
						</Li>
						<Li>
							<Checkbox>Nikon (12)</Checkbox>
						</Li>
						<Li>
							<Checkbox>BenQ (74)</Checkbox>
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
					Showing
					<span class="font-semibold text-gray-900 dark:text-white"
						>{currentPage * usersPerPage}-{currentPageUsers.length}</span
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
