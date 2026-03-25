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
			"¿Está activo?",
		],
		data: users
			.filter((user) => user.subrole === role)
			.map((user) => [
				user.id,
				user.name,
				user.lastname,
				user.email,
				user.document,
				user.phone,
				user.isActive ? "Sí" : "No",
			]),
	},

	columns: [{ select: 0, hidden: true }],
};
</script>

<Table dataTableOptions={dataTableOptions}/>
