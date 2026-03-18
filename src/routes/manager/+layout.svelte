<script lang="ts">
	let { children } = $props();

	import { Navbar, NavBrand, ButtonGroup, Button, Avatar } from 'flowbite-svelte';
	import { BellOutline } from 'flowbite-svelte-icons';
</script>

<Navbar>
	<NavBrand>
		<span>Administrador</span>
	</NavBrand>

	<ButtonGroup>
		<Button href="/manager/courses">Cursos</Button>
		<Button href="/manager/teachers">Profesores</Button>
		<Button href="/manager/students">Estudiantes</Button>
	</ButtonGroup>

	<Avatar class="size-7" />
</Navbar>

{@render children()}
