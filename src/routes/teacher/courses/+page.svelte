<script lang="ts">
	import { Table } from '@flowbite-svelte-plugins/datatable';
	import type { DataTable, DataTableOptions } from 'simple-datatables';

	const { data } = $props();

	const dataTableOptions = $derived({
		data: {
			headings: ['ID', 'Nombre', 'Descripción', 'Cupo'],
			data: data.courses.map((course) => [
				course.id,
				course.name,
				course.description,
				course.maxStudents
			])
		},
		columns: [{ select: 0, hidden: true }],
		paging: false,
		searchable: false
	}) satisfies DataTableOptions;
</script>

{#key data.courses}
	<Table {dataTableOptions} />
{/key}
