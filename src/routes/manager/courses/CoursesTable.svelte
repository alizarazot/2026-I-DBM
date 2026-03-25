<script lang="ts">
import type { DataTableOptions } from "@flowbite-svelte-plugins/datatable";
import { Table } from "@flowbite-svelte-plugins/datatable";

const { courses, users, onSelection } = $props();

const dataTableOptions = $derived({
	data: {
		headings: [
			"ID",
			"Nombre",
			"Descripción",
			"Día",
			"Hora",
			"Duración",
			"Cupo",
			"Profesor",
		],
		data: courses.map((course) => [
			course.id,
			course.name,
			course.description,
			course.day,
			course.startHour,
			course.duration,
			course.maxStudents,
			(() => {
				const teacher = users.find((user) => user.id === course.teacherId);
				if (!teacher) return "N/A";
				return `${teacher.name} ${teacher.lastName}`;
			})(),
		]),
	},
	columns: [{ select: 0, hidden: true }],
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
<Table dataTableOptions={dataTableOptions} selectable multiSelect={false} onSelectRow={idx => onSelection(dataTableOptions.data.data[idx][0])} />
{/key}
