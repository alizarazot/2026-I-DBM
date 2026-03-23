<script lang="ts">
import { Table } from "@flowbite-svelte-plugins/datatable";
import type { DataTableOptions } from "@flowbite-svelte-plugins/datatable";

import { users } from "$lib/dummy-db.json";

const dataTableOptions: DataTableOptions = {
	data: {
		headings: ["ID", "Nombre", "Correo electrónico", "Teléfono"],
		data: users
			.map(
				(user) =>
					(user.role === "teacher" && [
						user.id,
						user.fullName,
						user.email,
						user.phone,
					]) ||
					undefined,
			)
			.filter((item) => item),
	},
	columns: [{ select: 0, hidden: true }],
};
</script>

<Table dataTableOptions={dataTableOptions}/>
