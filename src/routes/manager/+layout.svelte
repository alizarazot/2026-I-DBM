<script lang="ts">
import { page } from "$app/state";

import { Navbar, NavBrand, ButtonGroup, Button, Avatar } from "flowbite-svelte";

const { children } = $props();

const sections = {
	courses: "Cursos",
	teachers: "Profesores",
	students: "Estudiantes",
};
</script>

<Navbar>
	<NavBrand>
		<span>Administrador</span>
	</NavBrand>

	<ButtonGroup>
		{#each Object.entries(sections) as section (section[0])}
				<Button href={"/manager/" + section[0]} class={{"text-primary-700": page.url.pathname.includes("/manager/" + section[0])}}> {section[1]} </Button>
		{/each}
	</ButtonGroup>

	<Avatar class="size-7" />
</Navbar>

{@render children()}
