<script lang="ts">
	let data = $state([]);
	setInterval(async () => {
		const req = await fetch('/broadcast/process-pcm');
		const json = await req.json();
		try {
			data = json;
		} catch (e) {}
	}, 1000);
</script>

{#each data as line}
	<p>{line[1]}</p>
{/each}
