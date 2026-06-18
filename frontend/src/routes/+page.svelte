<script lang="ts">
	import { Section, HeroHeader, HeroBody } from 'flowbite-svelte-blocks';
	import { Button } from 'flowbite-svelte';
	import { ArrowRightOutline } from 'flowbite-svelte-icons';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import type { User } from '$lib/api/model';

	onMount(async () => {
		// TODO: Show skeleton to avoid flash-of-content.
		try {
			const resp = await fetch('/api/auth');
			if (!resp.ok) {
				throw new Error('unknow error when asking for sign-in status');
			}

			const data = (await resp.json()) as { user: User };

			if (data.user.role == 'invalid') {
				throw new Error('user role is invalid');
			}

			// User is signed in, so redirect.
			goto(resolve(`/${data.user.role}`));
		} catch (err) {
			console.error(`Unexpected error: ${err}`);
		}
	});
</script>

<Section name="heroDefault">
	<HeroHeader>
		{#snippet h1()}
			<div class="mb-8 flex items-center justify-center gap-4">
				<img class="size-24" src="/images/logo.svg" alt="LÍNEA Logo" />
				<span>LÍNEA</span>
			</div>
		{/snippet}
		{#snippet h2()}Remote class management, <span class="underline">simplified</span>.{/snippet}
		{#snippet paragraph()}A unified platform for lecture administration, real-time attendance
			tracking, and automated student engagement.{/snippet}
	</HeroHeader>

	<div
		class="mb-8 flex flex-col space-y-4 sm:flex-row sm:justify-center sm:space-y-0 sm:space-x-4 lg:mb-16"
	>
		<a href={resolve('/sign-in')}>
			<Button size="lg" color="red">
				Sign in <ArrowRightOutline size="md" class="-mr-1 ml-2" />
			</Button>
		</a>
		<a href="https://github.com/alizarazot/2026-i-dbm">
			<Button size="lg" color="light">See the code on GitHub</Button>
		</a>
	</div>
	<HeroBody>
		{#snippet head()}Part of{/snippet}

		<div class="mt-8 flex flex-wrap items-center justify-center text-gray-500">
			<a
				title="UFPSO Logo"
				href="https://ufpso.edu.co/"
				class="mr-5 mb-5 hover:text-gray-800 lg:mb-0 dark:hover:text-gray-400"
			>
				<img class="w-64" src="/images/ufpso-logo.png" alt="UFPSO Logo" />
			</a>
		</div>
	</HeroBody>
</Section>
