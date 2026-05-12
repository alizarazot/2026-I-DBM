<script lang="ts">
	import type { PageProps } from './$types';
	import { Timeline, TimelineItem } from 'flowbite-svelte';
	import { onMount } from 'svelte';

	let { data }: PageProps = $props();

	onMount(() => {
		if (data.shouldAskLocation) {
			navigator.geolocation.getCurrentPosition(
				(position) => {
					fetch('/api/analytics/location', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({
							longitude: position.coords.longitude,
							latitude: position.coords.latitude
						})
					});
				},
				(error) => {
					console.error('Error getting location:', error);
				}
			);
		}
	});

	let lines = $state([]);

	let currentQuestion = $state<any>({});

	const updateLines = async () => {
		const req = await fetch('/broadcast/process-pcm');
		lines = await req.json();
		setTimeout(updateLines, 1000);
	};

	const updateQuestion = async () => {
		const resp = await fetch('/student/broadcast/questions');
		currentQuestion = await resp.json();
		setTimeout(updateQuestion, 1000);
	};
</script>

<div class="flex h-full flex-col">
	<Timeline class="m-8 grow overflow-y-scroll">
		{#each lines as line}
			<TimelineItem title={line} date=""><p></p></TimelineItem>
		{/each}
	</Timeline>

	{#if currentQuestion.question}
		<div>
			<h1>{currentQuestion.question}</h1>
			<ul>
				<li>- {currentQuestion.correctAnswer}</li>
				{#each currentQuestion.badAnswers as answer}
					<li>- {answer.answer}</li>
				{/each}
			</ul>
		</div>
	{/if}
</div>
