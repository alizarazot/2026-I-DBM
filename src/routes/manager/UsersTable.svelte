<script lang="ts">
import { Table } from "@flowbite-svelte-plugins/datatable";
import type { DataTableOptions } from "@flowbite-svelte-plugins/datatable";

const { role, users, onSelection } = $props();

const dataTableOptions = $derived({
	data: {
		headings: [
			"ID",
			"Nombres",
			"Apellidos",
			"Correo electrónico",
			"Documento",
			"Teléfono",
		],
		data: users
			.filter((user) => user.role === role)
			.map((user) => [
				user.id,
				user.name,
				user.lastName,
				user.email,
				user.document,
				user.phone,
			]),
	},
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

	columns: [{ select: 0, hidden: true }],
}) satisfies DataTableOptions;

$effect(() => {
	console.log(dataTableOptions);
});
</script>

{#key users}
<Table selectable onSelectRow={idx => onSelection(dataTableOptions.data.data[idx][0])} multiSelect={false} dataTableOptions={dataTableOptions}/>
{/key}
