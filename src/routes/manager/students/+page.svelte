<script lang="ts">
	import { Button } from 'flowbite-svelte';

	import Register from '../Register.svelte';
	import UsersTable from '../UsersTable.svelte';

	import type { PageData } from '../$types';
	import { invalidate } from '$app/navigation';

	let { data }: PageData = $props();

	let currentId = $state('');

	let registerOpenKind = $state<'register' | 'update' | null>(null);

	async function deleteUser() {
		const formData = new FormData();
		formData.append('id', currentId);
		const response = await fetch('/auth?/deleteUser', {
			method: 'POST',
			body: formData
		});
		invalidate('manager:users');
		currentId = '';
	}
</script>

<header class="me-1 flex justify-end gap-2 px-3">
	<Button
		onclick={() => {
			registerOpenKind = 'register';
		}}>Registrar estudiante</Button
	>
	<Button
		disabled={currentId === ''}
		onclick={() => {
			registerOpenKind = 'update';
		}}>Editar</Button
	>
	<Button disabled={currentId === ''} onclick={deleteUser}>Eliminar</Button>
</header>

<Register role="student" bind:openKind={registerOpenKind} bind:updateId={currentId} />

<div class="m-4 mt-0">
	<UsersTable
		role="student"
		users={data.users}
		onSelection={(id: string) => {
			currentId = id;
		}}
	/>
</div>
