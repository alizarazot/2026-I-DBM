<script lang="ts">
import { Button } from "flowbite-svelte";
import Register from "./Register.svelte";
import TeachersTable from "./TeachersTable.svelte";
let showRegister = $state(false);
</script>

<header class="flex justify-end px-3">
	<Button onclick={() => (showRegister = true)}>Registrar profesor</Button>
</header>

<Register bind:open={showRegister} />

<TeachersTable/>

