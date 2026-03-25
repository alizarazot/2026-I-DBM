<script lang="ts">
import { Table } from "@flowbite-svelte-plugins/datatable";
import type { DataTableOptions } from "@flowbite-svelte-plugins/datatable";

const { role, users } = $props();

const dataTableOptions: DataTableOptions = {
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

	columns: [{ select: 0, hidden: true }],
};
</script>

<Table dataTableOptions={dataTableOptions}/>
