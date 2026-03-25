<script lang="ts">
import { Button } from "flowbite-svelte";

import Register from "../Register.svelte";
//import UsersTable from "../UsersTable.svelte";

import { TODO } from "$lib/todo";

let registerOpenKind = $state<"register" | null>(null);
</script>

<header class="flex justify-end px-3">
	<Button onclick={() => (registerOpenKind = "register")}>Registrar profesor</Button>
	<Button onclick={TODO}>Editar</Button>
	<Button onclick={TODO}>Eliminar</Button>
</header>


<Register role="teachers" bind:openKind={registerOpenKind} />
