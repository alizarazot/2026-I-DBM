<script lang="ts">
import type { DataTableOptions } from "@flowbite-svelte-plugins/datatable";
import { Table } from "@flowbite-svelte-plugins/datatable";

import { courses } from "$lib/dummy-db.json";

const dataTableOptions: DataTableOptions = {
	data: {
		headings: ["ID", "Nombre", "Descripción", "Docente", "Fecha creación"],
		data: courses.map((course) => [
			course.id,
			course.title,
			course.description,
			course.teacherId,
			new Date(course.createdAt).toLocaleDateString(),
		]),
	},
	columns: [{ select: 0, hidden: true }],
};
</script>

<Table {dataTableOptions} />
