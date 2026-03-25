<script lang="ts">
import { Button } from "flowbite-svelte";
import { invalidate } from "$app/navigation";
import { downloadPdf } from "$lib/pdf";
import type { PageData } from "../$types";
import Register from "../Register.svelte";
import UsersTable from "../UsersTable.svelte";

let { data }: PageData = $props();

let currentId = $state("");

let registerOpenKind = $state<"register" | "update" | null>(null);

async function deleteUser() {
	const formData = new FormData();
	formData.append("id", currentId);
	const response = await fetch("/auth?/deleteUser", {
		method: "POST",
		body: formData,
	});
	invalidate("manager:users");
	currentId = "";
}

function downloadTablePdf() {
	const users = data.users.filter((user: any) => user.role === "student");
	downloadPdf(
		"Estudiantes",
		["Nombres", "Apellidos", "Correo electrónico", "Documento", "Teléfono"],
		users.map((u: any) => [u.name, u.lastName, u.email, u.document, u.phone]),
	);
}
</script>

<header class="flex justify-end px-3 gap-2 me-1">
	<Button onclick={() => {registerOpenKind = "register"}}>Registrar estudiante</Button>
	<Button disabled={currentId === ""} onclick={() => {registerOpenKind = "update"}}>Editar</Button>
	<Button disabled={currentId === ""} onclick={deleteUser}>Eliminar</Button>
	<Button onclick={downloadTablePdf}>PDF</Button>
</header>

<Register role="student" bind:openKind={registerOpenKind} bind:updateId={currentId} />

<div class="m-4 mt-0">
<UsersTable role="student" users={data.users} onSelection={(id: string)=>{currentId=id}}/>
</div>
