<script lang="ts">
	import { Section } from 'flowbite-svelte-blocks';
	import { Label, Input, Button, Select, Textarea } from 'flowbite-svelte';

	import type { CfcCategory } from '$lib/api/model';
	import { addCfc } from '$lib/api/common';

	const CATEGORIES: { value: CfcCategory; name: string }[] = [
		{ value: 'request', name: 'Request' },
		{ value: 'complaint', name: 'Complaint' },
		{ value: 'claim', name: 'Claim' },
		{ value: 'suggestion', name: 'Suggestion' }
	];

	let subject = $state('');
	let category = $state<CfcCategory | null>(null);
	let details = $state('');

	const handleSubmit = async (e: SubmitEvent): Promise<void> => {
		e.preventDefault();

		await addCfc(subject, category ?? 'invalid', details);

		subject = '';
		category = null;
		details = '';
	};
</script>

<Section name="crudcreateform">
	<h2 class="mb-4 text-xl font-bold text-gray-900 dark:text-white">
		Give feedback or make a complaint
	</h2>
	<form method="POST" onsubmit={handleSubmit}>
		<div class="grid grid-cols-1 gap-4 sm:gap-6 md:grid-cols-2">
			<div class="col-span-1">
				<Label for="subject" class="mb-2">Subject</Label>
				<Input name="subject" type="text" id="subject" required bind:value={subject} />
			</div>
			<div class="col-span-1">
				<Label>
					Category
					<Select class="mt-2" items={CATEGORIES} bind:value={category} name="category" required />
				</Label>
			</div>
			<div class="md:col-span-2">
				<Label for="details" class="mb-2">Details</Label>
				<Textarea
					id="details"
					class="w-full"
					rows={4}
					name="details"
					required
					bind:value={details}
				/>
			</div>
			<Button type="submit" class="w-32">Send {category}</Button>
		</div>
	</form>
</Section>
