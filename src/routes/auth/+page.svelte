<script lang="ts">
	import { Input, Label, ButtonGroup, RadioButton, Button } from 'flowbite-svelte';

	let userType = $state();
</script>

<div class="grid h-dvh w-dvw place-items-center">
	<form class="flex max-w-sm flex-col gap-3">
		<div>
			<Label>Correo electrónico</Label>
			<Input />
		</div>

		<div>
			<Label>Contraseña</Label>
			<Input />
		</div>

		<div class="mb-3 flex items-center justify-center gap-2">
			<Label>Rol:</Label>
			<ButtonGroup>
				<RadioButton value="student" bind:group={userType}>Estudiante</RadioButton>
				<RadioButton value="professor" bind:group={userType}>Profesor</RadioButton>
				<RadioButton value="manager" bind:group={userType}>Administrador</RadioButton>
			</ButtonGroup>
		</div>

		<Button>Iniciar sesión</Button>
	</form>
</div>
