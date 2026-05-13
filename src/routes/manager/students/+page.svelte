<script lang="ts">
	import { Button, Modal } from 'flowbite-svelte';
	import { Chart } from '@flowbite-svelte-plugins/chart';
	import type { ApexOptions } from 'apexcharts';

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
		await fetch('/auth?/deleteUser', {
			method: 'POST',
			body: formData
		});
		invalidate('manager:users');
		currentId = '';
	}

	const rawHours = new Set<number>([
		...Object.keys(data.maleAttendances ?? {}).map(Number),
		...Object.keys(data.femaleAttendances ?? {}).map(Number)
	]);
	const sortedHours = Array.from(new Set([...rawHours].flatMap((h) => [h - 1, h, h + 1]))).sort(
		(a, b) => a - b
	);
	const hoursLabels = sortedHours.map((h) => `${h}:00`);
	const maleAttendancesData = sortedHours.map((h) => data.maleAttendances?.[h] ?? 0);
	const femaleAttendancesData = sortedHours.map((h) => data.femaleAttendances?.[h] ?? 0);

	let showAttendanceModal = $state(false);
	let attendanceChartOptions: ApexOptions = {
		chart: {
			height: '400px',
			toolbar: {
				show: false
			}
		},
		series: [
			{
				name: 'Asistencias (hombres)',
				data: maleAttendancesData
			},
			{
				name: 'Asistencias (mujeres)',
				data: femaleAttendancesData,
				color: '#f6339a'
			}
		],
		xaxis: {
			categories: hoursLabels
		}
	};
</script>

<Modal bind:open={showAttendanceModal}>
	<Chart options={attendanceChartOptions} />
</Modal>

<header class="me-1 flex justify-end gap-2 px-3">
	<Button
		onclick={() => {
			showAttendanceModal = true;
		}}>Gráfica de asistencias</Button
	>
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
