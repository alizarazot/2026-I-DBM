<script lang="ts">
	import { uiHelpers, SidebarItem, SidebarGroup, Sidebar } from 'flowbite-svelte';
	import { ChartOutline, UsersOutline } from 'flowbite-svelte-icons';

	import Navbar from '$lib/components/Navbar.svelte';

	let { children } = $props();

	const sidebarState = uiHelpers();
</script>

<Navbar onClickSidebarButton={sidebarState.toggle} />

<Sidebar
	class="z-40 h-full"
	classes={{ div: 'pt-16' }}
	isOpen={sidebarState.isOpen}
	closeSidebar={sidebarState.close}
>
	<SidebarGroup>
		<SidebarItem label="Users" href="/manager/users">
			{#snippet icon()}
				<UsersOutline
					class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white"
				/>
			{/snippet}
		</SidebarItem>
		<SidebarItem label="Feedback & Complaints" href="/manager/feedback-and-complaints">
			{#snippet icon()}
				<ChartOutline
					class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white"
				/>
			{/snippet}
		</SidebarItem>
	</SidebarGroup>
</Sidebar>

<div class="md:ml-64">
	{@render children()}
</div>
