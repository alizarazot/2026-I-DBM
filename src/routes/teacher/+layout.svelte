<script lang="ts">
import {
	Navbar,
	NavBrand,
	ButtonGroup,
	Button,
	Avatar,
	Dropdown,
	DropdownDivider,
	DropdownHeader,
	DropdownItem,
} from "flowbite-svelte";

const { children, data } = $props();

async function signOut() {
	await fetch("/auth?/signOut", { method: "POST", body: new FormData() });
	window.location.href = "/auth";
}
</script>

<Navbar>
	<NavBrand>
		<span>Profesor</span>
	</NavBrand>

	<ButtonGroup>
		<Button href="/teacher/broadcast">Transmisión</Button>
		<Button href="/teacher/courses">Cursos</Button>
	</ButtonGroup>

	<div class="flex items-center gap-2">
		<span>{data.user.name}</span>
		<Avatar class="size-7" />
		<Dropdown simple placement="left">
	  	<DropdownHeader>
	  		<span class="block text-sm text-gray-900 dark:text-white">{data.user.name} {data.user.lastName}</span>
		    <span class="block truncate text-sm font-medium">{data.user.email}</span>
		  </DropdownHeader>
		  <DropdownDivider />
			<DropdownItem onclick={signOut}>Cerrar sesión </DropdownItem>
		</Dropdown>
	</div>
</Navbar>

{@render children()}
