<script lang="ts">
	import type { User } from '$lib/api/model';
	import { resolve } from '$app/paths';
	import {
		Avatar,
		Dropdown,
		DropdownGroup,
		DropdownHeader,
		DropdownItem,
		Navbar,
		NavBrand,
		SidebarButton
	} from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let { onClickSidebarButton } = $props();

	let user = $state<User>();
	onMount(async () => {
		const resp = await fetch('/api/auth');
		if (!resp.ok) {
			throw new Error('unknown error when asking for sign-in status');
		}

		const data = (await resp.json()) as { user: User };
		user = data.user;
	});

	const handleSignOut = async () => {
		const resp = await fetch('/api/sign-out');
		if (!resp.ok) {
			throw new Error('unknown error signing out');
		}

		goto(resolve('/'));
	};
</script>

<Navbar class="z-50 h-16 bg-gray-50">
	<SidebarButton onclick={onClickSidebarButton} />
	<NavBrand class="h-full gap-4" href="/">
		<img class="h-8" src="/images/logo.svg" alt="LÍNEA Logo" />
		<span class="capitalize">{user?.role}</span>
	</NavBrand>
	<div class="flex items-center md:order-2">
		<Avatar id="avatar-menu" />
	</div>
	<Dropdown placement="bottom" triggeredBy="#avatar-menu">
		<DropdownHeader>
			<span class="block text-sm"
				>{user?.info.firstName}
				{user?.info.middleName}
				{user?.info.firstSurname}
				{user?.info.secondSurname}</span
			>
			<span class="block truncate text-sm font-medium">{user?.email}</span>
		</DropdownHeader>
		<DropdownGroup>
			<DropdownItem onclick={handleSignOut} class="cursor-pointer">Sign out</DropdownItem>
		</DropdownGroup>
	</Dropdown>
</Navbar>
