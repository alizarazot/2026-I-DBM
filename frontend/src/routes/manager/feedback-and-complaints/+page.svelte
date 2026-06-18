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
		List,
		Li,
		Radio,
		Modal,
		P,
		Badge,
		Textarea
	} from 'flowbite-svelte';
	import { Section } from 'flowbite-svelte-blocks';
	import { FilterSolid } from 'flowbite-svelte-icons';

	import type { Cfc, CfcAnswer, CfcCategory } from '$lib/api/model';
	import { addCfcAnswer, getCfcAnswer, listCfcs } from '$lib/api/manager';

	let cfcs = $state<Cfc[]>([]);

	$effect(() => {
		listCfcs().then((c) => {
			cfcs = c;
		});
	});

	let searchTerm = $state('');
	let filterStatus = $state<'all' | 'pending' | 'answered'>('all');
	let filterCategory = $state<CfcCategory>('invalid');

	let filteredCfcs = $derived(
		cfcs.filter(
			(cfc) =>
				(filterStatus == 'pending' && !cfc.answered) ||
				(filterStatus == 'answered' && cfc.answered) ||
				(filterCategory != 'invalid' && cfc.category == filterCategory) ||
				(searchTerm != '' && cfc.subject.toLowerCase().indexOf(searchTerm) !== -1)
		)
	);

	let selectedCfc = $state<Cfc | null>(null);
	let cfcAnswer = $state<CfcAnswer | null>(null);
	let showModal = $state(false);
	const renderCfcs = $derived(
		searchTerm != '' || filterStatus != 'all' || filterCategory != 'invalid' ? filteredCfcs : cfcs
	);

	const handleShow = async (cfc: Cfc): Promise<void> => {
		selectedCfc = cfc;
		showModal = true;

		if (cfc.answered) {
			cfcAnswer = await getCfcAnswer(cfc.id);
			return;
		}

		cfcAnswer = null;
	};

	let formCfcAnswer = $state('');
	const handleCfcAnswer = async (e: SubmitEvent): Promise<void> => {
		e.preventDefault();
		showModal = false;

		await addCfcAnswer(selectedCfc?.id ?? '', formCfcAnswer);

		listCfcs().then((c) => {
			cfcs = c;
		});
	};
</script>

<Section name="advancedTable" class="bg-gray-50 p-3 sm:p-5 dark:bg-gray-900">
	<TableSearch
		placeholder="Search by subject"
		hoverable={true}
		classes={{
			root: 'bg-white dark:bg-gray-800 relative shadow-md sm:rounded-lg overflow-x-scroll',
			inner:
				'flex flex-col md:flex-row items-center justify-between space-y-3 md:space-y-0 md:space-x-4 p-4',
			search: 'w-full'
		}}
		bind:inputValue={searchTerm}
	>
		{#snippet header()}
			<div
				class="flex w-full flex-shrink-0 flex-col items-stretch justify-end space-y-2 md:w-auto md:flex-row md:items-center md:space-y-0 md:space-x-3"
			>
				<Button color="alternative">Status<FilterSolid class="ml-2 h-3 w-3 " /></Button>
				<Dropdown class="w-48 space-y-2 p-3 text-sm">
					<h6 class="mb-3 text-sm font-medium text-gray-900 dark:text-white">Choose a status</h6>
					<List tag="dl">
						<Li>
							<Radio value="all" bind:group={filterStatus}>All requests</Radio>
						</Li>
						<Li>
							<Radio value="pending" bind:group={filterStatus}>Pending</Radio>
						</Li>
						<Li>
							<Radio value="answered" bind:group={filterStatus}>Answered</Radio>
						</Li>
					</List>
				</Dropdown>
				<Button color="alternative">Category<FilterSolid class="ml-2 h-3 w-3 " /></Button>
				<Dropdown class="w-48 space-y-2 p-3 text-sm">
					<h6 class="mb-3 text-sm font-medium text-gray-900 dark:text-white">Choose a category</h6>
					<List tag="dl">
						<Li>
							<Radio value="invalid" bind:group={filterCategory}>All requests</Radio>
						</Li>
						<Li>
							<Radio value="request" bind:group={filterCategory}>Request</Radio>
						</Li>
						<Li>
							<Radio value="complaint" bind:group={filterCategory}>Complaint</Radio>
						</Li>
						<Li>
							<Radio value="claim" bind:group={filterCategory}>Claim</Radio>
						</Li>
						<Li>
							<Radio value="suggestion" bind:group={filterCategory}>Suggestion</Radio>
						</Li>
					</List>
				</Dropdown>
			</div>
		{/snippet}
		<TableHead>
			<TableHeadCell class="px-4 py-3" scope="col">Subject</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">Category</TableHeadCell>
			<TableHeadCell class="px-4 py-3" scope="col">Date</TableHeadCell>
			<TableHeadCell />
		</TableHead>
		<TableBody class="divide-y">
			{#each renderCfcs as cfc (cfc.id)}
				<TableBodyRow>
					<TableBodyCell class="px-4 py-3">{cfc.subject}</TableBodyCell>
					<TableBodyCell class="px-4 py-3 capitalize">{cfc.category}</TableBodyCell>
					<TableBodyCell class="px-4 py-3"
						>{cfc.updatedAt.toISOString().replace('T', ' ').slice(0, 16)}</TableBodyCell
					>
					<TableBodyCell class="px-4 py-3">
						{#if cfc.answered}
							<Button
								color="alternative"
								onclick={() => {
									handleShow(cfc);
								}}
							>
								View details</Button
							>
						{:else}
							<Button
								color="red"
								onclick={() => {
									handleShow(cfc);
								}}>Answer</Button
							>
						{/if}
					</TableBodyCell>
				</TableBodyRow>
			{/each}
		</TableBody>
	</TableSearch>
</Section>

<Modal title={selectedCfc?.subject} bind:open={showModal}>
	<P>
		<Badge class="capitalize">{selectedCfc?.category}</Badge>
		{selectedCfc?.userEmail}
	</P>
	<P class="pl-4">
		{selectedCfc?.details}
	</P>

	{#if cfcAnswer}
		<P>
			<Badge>{cfcAnswer.userEmail}</Badge>
			{cfcAnswer.answer}</P
		>
	{:else}
		<form method="POST" onsubmit={handleCfcAnswer}>
			<Textarea class="mb-4 w-full" bind:value={formCfcAnswer} required />
			<Button type="submit">Send answer</Button>
		</form>
	{/if}
</Modal>
