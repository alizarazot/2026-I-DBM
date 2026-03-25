<script lang="ts">
import { Button } from "flowbite-svelte";
import CoursesTable from "./CoursesTable.svelte";
import Register from "./Register.svelte";
import { invalidate } from "$app/navigation";

let currentId = $state("");

let registerOpenKind = $state<"register" | "update" | null>(null);

const { data } = $props();

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
</script>

<header class="flex justify-end px-3 gap-2 me-1">
	<Button onclick={() => {registerOpenKind = "register"}}>Añadir curso</Button>
	<Button disabled={currentId === ""} onclick={() => {registerOpenKind = "update"}}>Editar</Button>
	<Button disabled={currentId === ""} onclick={deleteCourse}>Eliminar</Button>
</header>

<Register users={data.users} bind:openKind={registerOpenKind} bind:updateId={currentId} />

<div class="m-4 mt-0">
<CoursesTable users={data.users} courses={data.courses} onSelection={(id: string)=>{currentId=id}}/>
</div>
