<script lang="ts">
	import { deserialize } from '$app/forms';
	import type { ActionResult } from '@sveltejs/kit';
	import { Label, Input, Button, Checkbox, ButtonGroup } from 'flowbite-svelte';
	import { MinusOutline, PlusOutline } from 'flowbite-svelte-icons';

	let searchQuery = $state('');
	let isQueryInverted = $state(false);
	let pageQueryNumber = $state(1);

	let searchResults = $state([]);

	const load = async () => {
		const bodyData = new FormData();
		bodyData.append('search-query', searchQuery);
		bodyData.append('is-query-inverted', isQueryInverted.toString());
		bodyData.append('page-query-number', (pageQueryNumber - 1).toString());

		const resp = await fetch('/manager/search?/searchQuery', { method: 'POST', body: bodyData });
		const result: ActionResult = deserialize(await resp.text());

		if (result.type != 'success' || !result.data) {
			throw new Error('Not successful request');
		}

		searchResults = result.data.searchResults;
	};
</script>

<div class="m-4 flex flex-col gap-4">
	<Label for="search-query">Busca cursos en la plataforma:</Label>
	<div class="flex gap-4">
		<Input id="search-query" clearable size="lg" bind:value={searchQuery} />
		<Checkbox bind:checked={isQueryInverted}>Invertir búsqueda</Checkbox>
		<ButtonGroup>
			<Button
				type="button"
				onclick={() => {
					pageQueryNumber = Math.max(1, pageQueryNumber - 1);
				}}><MinusOutline /></Button
			>
			<Input
				class="text-center"
				readonly
				pattern="[0-9]*"
				id="page-query-number"
				bind:value={pageQueryNumber}
			/>
			<Button
				type="button"
				onclick={() => {
					pageQueryNumber++;
				}}><PlusOutline /></Button
			>
		</ButtonGroup>
		<Button onclick={load}>Buscar</Button>
	</div>
</div>

<div class="m-4">
	{#if searchResults.length == 0 && searchQuery == ''}
		<span class="mt-32 block text-center text-gray-800"
			>Utiliza el campo de texto para iniciar una búsqueda.</span
		>
	{:else}
		<ul>
			{#each searchResults as result, idx (idx)}
				<li><strong>{result[0]}:</strong> {result[1]}</li>
			{/each}
		</ul>
	{/if}
</div>
