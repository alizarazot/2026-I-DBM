<script lang="ts">
	import { page } from '$app/state';
	import { redirect } from '@sveltejs/kit';

	import {
		Navbar,
		NavBrand,
		ButtonGroup,
		Button,
		Avatar,
		Dropdown,
		DropdownItem,
		DropdownHeader,
		DropdownDivider
	} from 'flowbite-svelte';

	const { children, data } = $props();

	const sections = {
		courses: 'Cursos',
		managers: 'Administradores',
		teachers: 'Profesores',
		students: 'Estudiantes'
	};

	async function signOut() {
		await fetch('/auth?/signOut', { method: 'POST', body: new FormData() });
		window.location.href = '/auth';
	}
</script>

<div class="flex h-dvh flex-col overflow-y-hidden">
	<div class="grow">
		<Navbar>
			<NavBrand>
				<span>Administrador</span>
			</NavBrand>

			<ButtonGroup>
				{#each Object.entries(sections) as section (section[0])}
					<Button
						href={'/manager/' + section[0]}
						class={{ 'text-primary-700': page.url.pathname.includes('/manager/' + section[0]) }}
					>
						{section[1]}
					</Button>
				{/each}
			</ButtonGroup>

			<div class="flex items-center gap-2">
				<span>{data.user.name}</span>
				<Avatar class="size-7" />
				<Dropdown simple placement="left">
					<DropdownHeader>
						<span class="block text-sm text-gray-900 dark:text-white"
							>{data.user.name} {data.user.lastName}</span
						>
						<span class="block truncate text-sm font-medium">{data.user.email}</span>
					</DropdownHeader>
					<DropdownDivider />
					<DropdownItem onclick={signOut}>Cerrar sesión</DropdownItem>
				</Dropdown>
			</div>
		</Navbar>
	</div>

	<div class="h-full grow">
		{@render children()}
	</div>
</div>
