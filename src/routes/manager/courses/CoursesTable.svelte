<script lang="ts">
import type { DataTableOptions } from "@flowbite-svelte-plugins/datatable";
import { Table } from "@flowbite-svelte-plugins/datatable";

const { courses } = $props();

console.log(courses);

const dataTableOptions: DataTableOptions = {
	data: {
		headings: [
			"Nombre",
			"Descripción",
			"Tipo",
			"Cantidad máxima de estudiantes",
		],
		data: courses.map((course) => [
			course.name,
			course.description,
			course.courseKind,
			course.maxStudents,
		]),
	},
};
</script>

<Table {dataTableOptions} />
