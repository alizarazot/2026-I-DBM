<script lang="ts">
import { enhance } from "$app/forms";
import type { ActionData } from "./$types";
import { Button, Input, Label, Heading, Span, Hr } from "flowbite-svelte";

let { form } = $props<{ form: ActionData }>();
</script>

<div class="grid h-dvh w-dvw place-items-center">
	<div class="flex flex-col items-center">
		<Heading class="text-center"> Plataforma <br> <Span highlight="orange">LINEA</Span> </Heading>
			<br>
	<form method=post action="?/signInEmail" class="flex w-full flex-col gap-3" use:enhance>
		<div>
			<Label>Correo electrónico</Label>
			<Input type=email name=email />
		</div>

		<div>
			<Label>Contraseña</Label>
			<Input  type="password" name=password />
		</div>

		<Button type=submit>Iniciar sesión</Button>
	</form>
	<p class="text-red-500"> {form?.message ?? ''} </p>
	</div>
</div>
