<script lang="ts">
	import type { PageProps } from './$types';
	import { Timeline, TimelineItem, GradientButton } from 'flowbite-svelte';
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
		const answer = await req.json();
		if (answer.lines.length > 0) {
			lines = answer.lines;
		}
		setTimeout(updateLines, 1000);
	};

	const updateQuestion = async () => {
		const resp = await fetch('/student/broadcast/questions');
		currentQuestion = await resp.json();
		setTimeout(updateQuestion, 1000);
	};

	onMount(() => {
		updateLines();
		updateQuestion();
	});

	const buttonColors = ['blue', 'green', 'cyan', 'teal'];
	const chooseRandomColor = (): 'blue' | 'green' | 'cyan' | 'teal' =>
		buttonColors[Math.floor(Math.random() * buttonColors.length)] as
			| 'blue'
			| 'green'
			| 'cyan'
			| 'teal';

	async function sendAnswer(answer: number) {
		const formData = new FormData();
		formData.append('answer', answer.toString());

		await fetch('/student/broadcast/questions', {
			method: 'POST',
			body: formData
		});
	}
</script>

<div class="flex h-full flex-col">
	<Timeline class="m-8 grow snap-y snap-proximity overflow-y-scroll">
		{#each lines as line, key (key)}
			<TimelineItem class="snap-align-none" title={line} date=""><p></p></TimelineItem>
		{/each}
		<TimelineItem class="snap-end" title="" date=""><p></p></TimelineItem>
	</Timeline>

	{#if currentQuestion.question}
		<div class="border-t pt-3">
			<h1 class="mx-4"><strong>Pregunta:</strong> {currentQuestion.question}</h1>
			<ul class=" flex gap-4 p-4">
				<GradientButton color={chooseRandomColor()} onclick={() => sendAnswer(0)}>
					{currentQuestion.correctAnswer}
				</GradientButton>
				{#each currentQuestion.badAnswers as answer, id (id)}
					<GradientButton color={chooseRandomColor()} onclick={() => sendAnswer(id + 1)}>
						{answer.answer}
					</GradientButton>
				{/each}
			</ul>
		</div>
	{/if}
</div>
