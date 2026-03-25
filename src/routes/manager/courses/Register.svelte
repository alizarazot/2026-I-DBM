<script lang="ts">
import { deserialize, enhance } from "$app/forms";
import { Button, Input, Label, Modal, Select } from "flowbite-svelte";

let { openKind = $bindable(null), updateId = $bindable(""), users } = $props();

let open = $state(false);

$effect(() => {
	open = openKind != null;
});

let action = $derived.by(() => {
	switch (openKind) {
		case null:
			return "";
		case "register":
			return "/manager/courses?/addCourse";
		case "update":
			return "/manager/courses?/editCourse";
	}
});

let actionLabel = $derived.by(() => {
	switch (openKind) {
		case null:
			return "";
		case "register":
			return "Añadir";
		case "update":
			return "Editar";
	}
});

let name = $state("");
let description = $state("");
let maxStudents = $state("");
let teacherId = $state("");

$effect(() => {
	const id = updateId;
	if (id === "") return;

	(async function () {
		const formData = new FormData();
		formData.append("id", id);
		const response = await fetch("/manager/courses?/getCourse", {
			method: "POST",
			body: formData,
		});
		const text = await response.text();
		const result = deserialize(text);
		if (result.type != "success") {
			return;
		}
		const course = result.data!.course as any;
		name = course.name;
		description = course.description;
		maxStudents = course.maxStudents;
		teacherId = course.teacherId;
	})();
});
const teachers = $derived.by(() =>
	users
		.filter((user) => user.role === "teacher")
		.map((user) => ({
			value: user.id,
			name: `${user.name} ${user.lastName}`,
		})),
);

let errorMsg = $state();
</script>

<Modal open={open} onclose={()=>{openKind = null}} size="xs">
		<form method=POST action={action} use:enhance class="flex flex-col space-y-6">
		<h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">{actionLabel} curso</h3>

		{#if openKind === "update"}
		<Input name=id type=hidden bind:value={updateId}/>
		{/if}

		<Label class="space-y-2">
			<span>Nombre</span>
			<Input name=name maxlength={100} required bind:value={name} />
		</Label>

		<Label class="space-y-2">
			<span>Descripción</span>
			<Input name=description maxlength={200} required bind:value={description}/>
		</Label>

		<Label class="space-y-2">
			<span>Cantidad máxima de estudiantes</span>
			<Input name=maxStudents type=number required bind:value={maxStudents}/>
		</Label>
		
		<Label class="space-y-2">
			<span>Profesor</span>
			<Select name=teacherId items={teachers} placeholder="Seleccionar profesor..." bind:value={teacherId}/>
		</Label>

		<div class="flex justify-end gap-2">
			<Button type=submit>{actionLabel}</Button>
			<Button onclick={()=> {openKind =null}}>Cancelar</Button>
		</div>
		<p class="text-red-500">{errorMsg}</p>
	</form>
</Modal>
