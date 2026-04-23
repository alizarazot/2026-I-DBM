<script lang="ts">
	import { Timeline, TimelineItem, Button } from 'flowbite-svelte';
	import { onMount } from 'svelte';

	let lines = $state([]);

	const updateLines = async () => {
		const req = await fetch('/broadcast/process-pcm');
		lines = await req.json();
		setTimeout(updateLines, 1000);
	};

	onMount(() => {
		updateLines();
	});
</script>

<div class="flex h-full flex-col">
	<Timeline class="m-8 grow overflow-y-scroll">
		{#each lines as line}
			<TimelineItem title={line} date=""><p></p></TimelineItem>
		{/each}
	</Timeline>

	<div></div>
</div>
