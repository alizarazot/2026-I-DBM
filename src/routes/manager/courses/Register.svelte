<script lang="ts">
import { enhance } from "$app/forms";
import { Button, Input, Label, Modal, Select } from "flowbite-svelte";

let { openKind = $bindable(false), updateId = $bindable(""), users } = $props();

const teachers = $derived.by(() =>
	users
		.filter((user) => user.role === "teacher")
		.map((user) => ({
			value: user.id,
			name: `${user.name} ${user.lastName}`,
		})),
);
</script>

<Modal open={openKind != null} size="xs">
		<form method=POST action="/manager/courses?/addCourse" use:enhance class="flex flex-col space-y-6">
		<h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Añadir curso</h3>

		<Label class="space-y-2">
			<span>Nombre</span>
			<Input name=name maxlength={100} required />
		</Label>

		<Label class="space-y-2">
			<span>Descripción</span>
			<Input name=description maxlength={200} required/>
		</Label>

		<Label class="space-y-2">
			<span>Cantidad máxima de estudiantes</span>
			<Input name=maxStudents type=number required />
		</Label>
		
		<Label class="space-y-2">
			<span>Profesor</span>
			<Select name=teacherId items={teachers} placeholder="Seleccionar profesor..."/>
		</Label>

		<div class="flex justify-end gap-2">
			<Button type=submit>Añadir</Button>
			<Button onclick={close}>Cancelar</Button>
		</div>
	</form>
</Modal>
