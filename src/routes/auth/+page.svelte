<script lang="ts">
import { enhance } from "$app/forms";
import type { ActionData } from "./$types";
import {
	Button,
	ButtonGroup,
	Input,
	Label,
	RadioButton,
} from "flowbite-svelte";

let { form } = $props<{ form: ActionData }>();
</script>

<div class="grid h-dvh w-dvw place-items-center">
	<div>
	<form method=post action="?/signInEmail" class="flex max-w-sm flex-col gap-3" use:enhance>
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
