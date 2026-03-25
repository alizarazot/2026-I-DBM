<script lang="ts">
import type { DataTableOptions } from "@flowbite-svelte-plugins/datatable";
import { Table } from "@flowbite-svelte-plugins/datatable";

const { courses, users } = $props();

console.log(users, courses);

const dataTableOptions = $derived({
	data: {
		headings: ["Nombre", "Descripción", "Cupo", "Profesor"],
		data: courses.map((course) => [
			course.name,
			course.description,
			course.maxStudents,
			(() => {
				const teacher = users.find((user) => user.id === course.teacherId);
				if (!teacher) return "N/A";
				return `${teacher.name} ${teacher.lastName}`;
			})(),
		]),
	},
	paging: false,
	searchable: false,
	rowRender: (row: any, tr: any, _index: number) => {
		if (!tr.attributes) {
			tr.attributes = {};
		}
		if (!tr.attributes.class) {
			tr.attributes.class = "";
		}
		if (row.selected) {
			tr.attributes.class += " selected";
		} else {
			tr.attributes.class = tr.attributes.class.replace(" selected", "");
		}
		return tr;
	},
}) satisfies DataTableOptions;
</script>

{#key courses}
<Table dataTableOptions={dataTableOptions} selectable multiSelect={false} />
{/key}
