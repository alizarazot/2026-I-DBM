<script lang="ts">
import { Button } from "flowbite-svelte";

import Register from "../Register.svelte";
import UsersTable from "../UsersTable.svelte";

import type { PageData } from "../$types";

import { TODO } from "$lib/todo";

let { data }: PageData = $props();

let currentId = $state("");

let registerOpenKind = $state<"register" | "update" | null>(null);
</script>

<header class="flex justify-end px-3">
	<Button onclick={() => {registerOpenKind = "register"}}>Registrar administrador</Button>
	<Button disabled={currentId === ""} onclick={() => {registerOpenKind = "update"}}>Editar</Button>
	<Button onclick={TODO}>Eliminar</Button>
</header>


<Register role="manager" bind:openKind={registerOpenKind} bind:updateId={currentId} />
<UsersTable role="manager" users={data.users} onSelection={(id: string)=>{currentId=id}}/>
