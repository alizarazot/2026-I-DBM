<script lang="ts">
let { children } = $props();

import { Navbar, NavBrand, ButtonGroup, Button, Avatar } from "flowbite-svelte";
import { BellOutline } from "flowbite-svelte-icons";
</script>

<Navbar>
	<NavBrand>
		<span>VoiceClass</span>
	</NavBrand>

	<ButtonGroup>
		<Button href="/student/broadcast">Transmisión</Button>
		<Button href="/student/tasks">Tareas</Button>
		<Button href="/student/details">Detalles</Button>
	</ButtonGroup>

	<div class="flex items-center gap-2">
		<BellOutline class="size-7" />
		<Avatar class="size-7" />
	</div>
</Navbar>

{@render children()}
