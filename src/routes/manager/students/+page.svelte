<script lang="ts">
import { Button } from "flowbite-svelte";

import Register from "../Register.svelte";
import UsersTable from "../UsersTable.svelte";

import type { PageData } from "../$types";

import { TODO } from "$lib/todo";

let { data }: PageData = $props();

let registerOpenKind = $state<"register" | null>(null);
</script>

<header class="flex justify-end px-3">
	<Button onclick={() => (registerOpenKind = "register")}>Registrar estudiante</Button>
	<Button onclick={TODO}>Editar</Button>
	<Button onclick={TODO}>Eliminar</Button>
</header>


<Register role="student" bind:openKind={registerOpenKind} />
<UsersTable role="student" users={data.users} />
