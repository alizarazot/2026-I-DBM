<script lang="ts">
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

	let { onClickSidebarButton } = $props();
</script>

<Navbar class="z-50 h-16 bg-gray-50">
	<SidebarButton onclick={onClickSidebarButton} />
	<NavBrand class="h-full gap-4" href="/">
		<img class="h-8" src="/images/logo.svg" alt="LÍNEA Logo" />
		<span>Teacher</span>
	</NavBrand>
	<div class="flex items-center md:order-2">
		<Avatar id="avatar-menu" />
	</div>
	<Dropdown placement="bottom" triggeredBy="#avatar-menu">
		<DropdownHeader>
			<span class="block text-sm">Bonnie Green</span>
			<span class="block truncate text-sm font-medium">name@flowbite.com</span>
		</DropdownHeader>
		<DropdownGroup>
			<DropdownItem class="cursor-pointer">Sign out</DropdownItem>
		</DropdownGroup>
	</Dropdown>
</Navbar>
