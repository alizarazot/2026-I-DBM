<script lang="ts">
import { page } from "$app/state";

import { Navbar, NavBrand, ButtonGroup, Button, Avatar } from "flowbite-svelte";

const { children, data } = $props();

const sections = {
	courses: "Cursos",
	managers: "Administradores",
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

	<div class="flex items-center gap-2">
		<span>{data.user.name}</span>
	<Avatar class="size-7" />
	</div>
</Navbar>

{@render children()}
