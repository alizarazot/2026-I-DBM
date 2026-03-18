<script lang="ts">
	import { goto } from '$app/navigation';
	import { Input, Label, ButtonGroup, RadioButton, Button } from 'flowbite-svelte';
	import { userState } from '$lib/user.svelte.ts';

	let email = $state();
	let password = $state();
	let userType = $state('student');

	function doLogin() {
		console.log('login called');
		userState.userType = userType;
		goto('/');
	}
</script>

<div class="grid h-dvh w-dvw place-items-center">
	<form on:submit|preventDefault={doLogin} class="flex max-w-sm flex-col gap-3">
		<div>
			<Label>Correo electrónico</Label>
			<Input bind:value={email} />
		</div>

		<div>
			<Label>Contraseña</Label>
			<Input bind:value={password} type="password" />
		</div>

		<div class="mb-3 flex items-center justify-center gap-2">
			<Label>Rol:</Label>
			<ButtonGroup>
				<RadioButton value="student" bind:group={userType}>Estudiante</RadioButton>
				<RadioButton value="professor" bind:group={userType}>Profesor</RadioButton>
				<RadioButton value="manager" bind:group={userType}>Administrador</RadioButton>
			</ButtonGroup>
		</div>

		<Button type="submit">Iniciar sesión</Button>
	</form>
</div>
