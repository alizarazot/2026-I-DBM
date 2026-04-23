<script lang="ts">
	import { Timeline, TimelineItem, Button } from 'flowbite-svelte';
	import { onMount } from 'svelte';

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

	onMount(() => {
		updateLines();
		updateQuestion();
	});
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
