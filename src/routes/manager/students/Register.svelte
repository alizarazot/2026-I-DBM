<script lang="ts">
import { Button, Input, Label, Modal, Radio } from "flowbite-svelte";

let { open = $bindable(false) } = $props();

let isActive = $state("true");

function close() {
	open = false;
}
</script>

<Modal form bind:open size="xs">
	<div class="flex flex-col space-y-6">
		<h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Registrar estudiante</h3>

		<Label class="space-y-2">
			<span>Número de cédula de ciudadanía</span>
			<Input type="number" />
		</Label>

		<div class="flex gap-2 space-y-2">
			<Label class="grow">
				<span>Nombre</span>
				<Input />
			</Label>

			<Label class="grow">
				<span>Apellido</span>
				<Input />
			</Label>
		</div>

		<div class="flex gap-2 space-y-2">
			<Label class="grow">
				<span>Correo electrónico</span>
				<Input type="email" />
			</Label>

			<Label class="grow">
				<span>Teléfono</span>
				<Input type="tel" />
			</Label>
		</div>

		<Label class="space-y-2">
			<span>Fecha de nacimiento</span>
			<Input type="date" />
		</Label>

		<Label class="space-y-2">
			<span>Contraseña inicial</span>
			<Input type="password" />
		</Label>

		<Label class="flex items-center gap-2 space-y-2">
			<span class="m-0">Estado</span>
			<ul
				class="flex w-full items-center divide-x divide-gray-200 rounded-lg border border-gray-200"
			>
				<li class="w-full">
					<Radio classes={{ label: 'p-3' }} name="is-active" value="true" bind:group={isActive}
						>Activo</Radio
					>
				</li>
				<li class="w-full">
					<Radio classes={{ label: 'p-3' }} name="is-active" value="false" bind:group={isActive}
						>Inactivo</Radio
					>
				</li>
			</ul>
		</Label>

		<div class="flex justify-end gap-2">
			<Button onclick={close}>Cancelar</Button>
			<Button onclick={close}>Guardar</Button>
		</div>
	</div>
</Modal>
