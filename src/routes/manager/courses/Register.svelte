<script lang="ts">
import { Button, Input, Label, Modal } from "flowbite-svelte";

let { open = $bindable(false) } = $props();

function close() {
	open = false;
}
</script>

<Modal bind:open size="xs">
		<form method=POST action="/manager/courses?/addCourse" class="flex flex-col space-y-6">
		<h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Registrar curso</h3>

		<Label class="space-y-2">
			<span>Nombre del curso</span>
			<Input name=name maxlength={100} />
		</Label>

		<Label class="space-y-2">
			<span>Descripción</span>
			<Input name=description maxlength={200} />
		</Label>
		
		<Label class="space-y-2">
			<span>Tipo del curso</span>
			<Input name=courseKind maxlength={200} />
		</Label>

		<Label class="space-y-2">
			<span>Profesor</span>
			<Input />
		</Label>

		
		<Label class="space-y-2">
			<span>Cantidad máxima de estudiantes</span>
			<Input name=maxStudents type=number />
		</Label>

		<div class="flex justify-end gap-2">
			<Button type=submit>Registrar</Button>
			<Button onclick={close}>Cancelar</Button>
		</div>
	</form>
</Modal>
