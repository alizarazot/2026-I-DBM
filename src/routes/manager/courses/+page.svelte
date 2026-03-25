<script lang="ts">
import { Button } from "flowbite-svelte";
import { invalidate } from "$app/navigation";
import { downloadPdf } from "$lib/pdf";
import CoursesTable from "./CoursesTable.svelte";
import Register from "./Register.svelte";

let currentId = $state("");

let registerOpenKind = $state<"register" | "update" | null>(null);

let { data } = $props();

async function deleteCourse() {
	const formData = new FormData();
	formData.append("id", currentId);
	await fetch("/manager/courses?/deleteCourse", {
		method: "POST",
		body: formData,
	});
	invalidate("manager:courses");
	currentId = "";
}

function downloadTablePdf() {
	downloadPdf(
		"Cursos",
		["Nombre", "Descripción", "Cupo", "Profesor"],
		data.courses.map((c: any) => [
			c.name,
			c.description,
			c.maxStudents,
			data.users.find((u: any) => u.id === c.teacherId)?.name +
				" " +
				data.users.find((u: any) => u.id === c.teacherId)?.lastName || "N/A",
		]),
	);
}
</script>

<div class="flex flex-col h-full overflow-y-hidden">
<header class="flex justify-end px-3 gap-2 me-1">
	<Button onclick={() => {registerOpenKind = "register"}}>Añadir curso</Button>
	<Button disabled={currentId === ""} onclick={() => {registerOpenKind = "update"}}>Editar</Button>
	<Button disabled={currentId === ""} onclick={deleteCourse}>Eliminar</Button>
	<Button onclick={downloadTablePdf}>PDF</Button>
</header>

<Register users={data.users} bind:openKind={registerOpenKind} bind:updateId={currentId} />

<div class="mx-4 pb-4 overflow-y-auto grow h-full">
<CoursesTable users={data.users} courses={data.courses} onSelection={(id: string)=>{currentId=id}}/>
</div>
</div>
