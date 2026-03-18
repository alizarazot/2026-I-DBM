<script lang="ts">
	import { Button } from 'flowbite-svelte';
	import Register from './Register.svelte';

	let showRegister = $state(false);
</script>

<header class="flex justify-end px-3">
	<Button onclick={(_) => (showRegister = true)}>Registrar profesor</Button>
</header>

<Register bind:open={showRegister} />
