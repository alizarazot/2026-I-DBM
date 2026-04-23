<script lang="ts">
	import { Button } from 'flowbite-svelte';
	import CoursesTable from './CoursesTable.svelte';
	import Register from './Register.svelte';
	import { invalidate } from '$app/navigation';

	let currentId = $state('');

	let registerOpenKind = $state<'register' | 'update' | null>(null);

	const { data } = $props();

	async function deleteCourse() {
		const formData = new FormData();
		formData.append('id', currentId);
		await fetch('/manager/courses?/deleteCourse', {
			method: 'POST',
			body: formData
		});
		invalidate('manager:courses');
		currentId = '';
	}
</script>

<div class="flex h-full flex-col overflow-y-hidden">
	<header class="me-1 flex justify-end gap-2 px-3">
		<Button
			onclick={() => {
				registerOpenKind = 'register';
			}}>Añadir curso</Button
		>
		<Button
			disabled={currentId === ''}
			onclick={() => {
				registerOpenKind = 'update';
			}}>Editar</Button
		>
		<Button disabled={currentId === ''} onclick={deleteCourse}>Eliminar</Button>
	</header>

	<Register users={data.users} bind:openKind={registerOpenKind} bind:updateId={currentId} />

	<div class="mx-4 h-full grow overflow-y-auto pb-4">
		<CoursesTable
			users={data.users}
			courses={data.courses}
			onSelection={(id: string) => {
				currentId = id;
			}}
		/>
	</div>
</div>
