<script lang="ts">
import { Button } from "flowbite-svelte";
import Register from "./Register.svelte";
import StudentsTable from "./StudentsTable.svelte";

let showRegister = $state(false);
</script>

<header class="flex justify-end px-3">
	<Button onclick={() => (showRegister = true)}>Registrar estudiante</Button>
</header>

<Register bind:open={showRegister} />

<StudentsTable />
