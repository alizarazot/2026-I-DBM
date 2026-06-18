<script lang="ts">
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { Section, Register } from 'flowbite-svelte-blocks';
	import { Button, Label, Input } from 'flowbite-svelte';

	onMount(async () => {
		// TODO: Show skeleton to avoid flash-of-content.
		try {
			const resp = await fetch('/api/auth');

			if (resp.ok) {
				// User is signed in, so redirect.
				goto(resolve('/'));
			}
		} catch (err) {
			console.error(`Unexpected error: ${err}`);
		}
	});

	const handleSignIn = async (e: SubmitEvent) => {
		e.preventDefault();

		const target = e.target as HTMLFormElement;

		const formData = new FormData(target);
		const data = Object.fromEntries(formData);

		try {
			const resp = await fetch(target.action, {
				method: target.method,
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			});

			if (resp.status == 401) {
				throw new Error('Invalid credentials');
			}

			if (!resp.ok) {
				throw new Error(`Response status: '${resp.status}'`);
			}

			// At this point, the server is in charge of managing the auth cookie (only applies to `/api/*`).

			goto(resolve('/'));
		} catch (err) {
			// TODO: Show dialog with error message.
			console.error(err);
		}
	};
</script>

<Section class="grid h-screen" name="login">
	<Register href="/">
		{#snippet top()}
			<img class="mr-2 h-8 w-8" src="/images/logo.svg" alt="logo" />
			LÍNEA
		{/snippet}
		<div class="space-y-4 p-6 sm:p-8 md:space-y-6">
			<form
				class="flex flex-col space-y-6"
				method="POST"
				action="/api/auth"
				onsubmit={handleSignIn}
			>
				<h3 class="p-0 text-xl font-medium text-gray-900 dark:text-white">Welcome back!</h3>
				<Label class="space-y-2">
					<span>Your email</span>
					<Input type="email" name="email" required />
				</Label>
				<Label class="space-y-2">
					<span>Your password</span>
					<Input type="password" name="password" required />
				</Label>

				<!-- TODO: Implement in backend. -->
				<!--div class="flex items-start">
					<Checkbox>Remember me</Checkbox>
					<a href="/" class="ml-auto text-sm text-blue-700 hover:underline dark:text-blue-500"
						>Forgot password?</a
					>
				</div-->
				<Button type="submit" class="w-full">Sign in</Button>
			</form>
		</div>
	</Register>
</Section>
