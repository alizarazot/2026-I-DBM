<script lang="ts">
import { Button, Input, Label, Modal, Radio } from "flowbite-svelte";

import { deserialize, enhance } from "$app/forms";
import { invalidate } from "$app/navigation";

let { openKind = $bindable(null), role, updateId = $bindable("") } = $props();

let open = $state(false);

$effect(() => {
	open = openKind != null;
});

let action = $derived.by(() => {
	switch (openKind) {
		case null:
			return "";
		case "register":
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
		case "update":
			return "/auth?/updateUser";
		default:
			throw Error("Invalid openKind");
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

let document = $state("");
let firstName = $state("");
let lastName = $state("");
let email = $state("");
let password = $state("");
let phone = $state("");

$effect(() => {
	const id = updateId;
	if (id === "") {
		return;
	}

	(async function () {
		const formData = new FormData();
		formData.append("id", id);
		const response = await fetch("/manager?/getUser", {
			method: "POST",
			body: formData,
		});
		const text = await response.text();
		const result = deserialize(text);
		if (result.type != "success") {
			return;
		}
		const user = result.data!.user as any;
		document = user.document;
		firstName = user.name;
		lastName = user.lastName;
		email = user.email;
		phone = user.phone;
	})();
});

$effect(() => {
	if (openKind === "register") {
		document = "";
		firstName = "";
		lastName = "";
		email = "";
		phone = "";
	}
});

const registerButtonText = $derived.by(() => {
	switch (openKind) {
		case null:
			return "";
		case "register":
			return "Registrar";
		case "update":
			return "Editar";
	}
});

let errorMsg = $state();
</script>

<Modal open={open} onclose={() => {openKind=null}}  size="xs">
		<form method=POST action={action} use:enhance={()=>{ return async ({ result, update }) => {
        await invalidate("manager:users");
        await update();
        openKind = null
    };}}
    class="flex flex-col space-y-6">

		<h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">{registerButtonText} {registerLabel}</h3>
		
		{#if openKind === "update"}
		<Input name=id type=hidden bind:value={updateId}/>
		{/if}

		<Label class="space-y-2">
			<span>Número de cédula de ciudadanía</span>
			<Input name=document type="number" required bind:value={document}/>
		</Label>

		<div class="flex gap-2 space-y-2">
			<Label class="grow">
				<span>Nombres</span>
				<Input name=firstName required bind:value={firstName}/>
			</Label>

			<Label class="grow">
				<span>Apellidos</span>
				<Input name=lastName required bind:value={lastName}/>
			</Label>
		</div>

		<div class="flex gap-2 space-y-2">
			<Label class="grow">
				<span>Correo electrónico</span>
				<Input name=email type="email" required bind:value={email}/>
			</Label>

			<Label class="grow">
				<span>Teléfono</span>
				<Input name=phone type="tel" required bind:value={phone} />
			</Label>
		</div>

		{#if false}
		<Label class="space-y-2">
			<span>Especialidad</span>
			<Input name=speciality required />
		</Label>
		{/if}

		{#if openKind === "register"}
		<Label class="space-y-2">
			<span>Contraseña inicial</span>
			<Input name=password type="password" required bind:value={password}/>
		</Label>
		{/if}

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
			<Button type="submit">{registerButtonText}</Button>
			<Button onclick={() => {openKind = null}}>Cancelar</Button>
		</div>
		<p class="text-red-500">{errorMsg}</p>
	</form>
</Modal>
