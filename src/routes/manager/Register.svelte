<script lang="ts">
import { Button, Input, Label, Modal, Radio } from "flowbite-svelte";

import { enhance } from "$app/forms";

let { openKind = $bindable(null), role } = $props();
let open = $derived(openKind != null);
let action = $derived.by(() => {
	switch (role) {
		case "manager":
			return "/auth?/addManager";
		case "teacher":
			return "/auth?/addTeacher";
		case "student":
			return "/auth?/addStudent";
		default:
			throw Error("Invalid role");
	}
});
let registerLabel = $derived.by(() => {
	switch (role) {
		case "manager":
			return "administrador";
		case "teacher":
			return "profesor";
		case "student":
			return "estudiante";
		default:
			throw Error("Invalid role");
	}
});

let errorMsg = $state();
</script>

<Modal bind:open={open} size="xs">
		<form method=POST action={action} use:enhance class="flex flex-col space-y-6">
		<h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Registrar {registerLabel}</h3>

		<Label class="space-y-2">
			<span>Número de cédula de ciudadanía</span>
			<Input name=document type="number" required/>
		</Label>

		<div class="flex gap-2 space-y-2">
			<Label class="grow">
				<span>Nombres</span>
				<Input name=firstName required/>
			</Label>

			<Label class="grow">
				<span>Apellidos</span>
				<Input name=lastName required />
			</Label>
		</div>

		<div class="flex gap-2 space-y-2">
			<Label class="grow">
				<span>Correo electrónico</span>
				<Input name=email type="email" required />
			</Label>

			<Label class="grow">
				<span>Teléfono</span>
				<Input name=phone type="tel" required />
			</Label>
		</div>

		{#if false}
		<Label class="space-y-2">
			<span>Especialidad</span>
			<Input name=speciality required />
		</Label>
		{/if}

		<Label class="space-y-2">
			<span>Contraseña inicial</span>
			<Input name=password type="password" required/>
		</Label>

		{#if false}
		<Label class="flex items-center gap-2 space-y-2">
			<span class="m-0">Estado</span>
			<ul
				class="flex w-full items-center divide-x divide-gray-200 rounded-lg border border-gray-200"
			>
				<li class="w-full">
					<Radio classes={{ label: 'p-3' }} name="is-active" value="true" 
						>Activo</Radio
					>
				</li>
				<li class="w-full">
					<Radio classes={{ label: 'p-3' }} name="is-active" value="false"
						>Inactivo</Radio
					>
				</li>
			</ul>
		</Label>
		{/if}

		<div class="flex justify-end gap-2">
			<Button onclick={() => {openKind = null}}>Cancelar</Button>
			<Button type="submit">Registrar</Button>
		</div>
		<p class="text-red-500">{errorMsg}</p>
	</form>
</Modal>
