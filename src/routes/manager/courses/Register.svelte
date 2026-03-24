<script lang="ts">
import { Button, Input, Label, Modal } from "flowbite-svelte";

let { open = $bindable(false) } = $props();

function close() {
	open = false;
}
</script>

<Modal form bind:open size="xs">
	<div class="flex flex-col space-y-6">
		<h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Registrar curso</h3>

		<Label class="space-y-2">
			<span>Nombre del curso</span>
			<Input maxlength={100} />
		</Label>

		<Label class="space-y-2">
			<span>Descripción</span>
			<Input maxlength={200} />
		</Label>

		<Label class="space-y-2">
			<span>Documento del docente</span>
			<Input />
		</Label>

		<div class="flex justify-end gap-2">
			<Button onclick={close}>Cancelar</Button>
			<Button onclick={close}>Guardar</Button>
		</div>
	</div>
</Modal>
