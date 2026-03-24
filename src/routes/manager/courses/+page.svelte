<script lang="ts">
import { Button } from "flowbite-svelte";
import CoursesTable from "./CoursesTable.svelte";
import Register from "./Register.svelte";

let showRegister = $state(false);
</script>

<header class="flex justify-end px-3">
	<Button onclick={(_: any) => (showRegister = true)}>Registrar curso</Button>
</header>

<Register bind:open={showRegister} />

<CoursesTable />
