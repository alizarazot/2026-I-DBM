<script lang="ts">
	import { SvelteMap } from 'svelte/reactivity';
	import { onMount } from 'svelte';
	import { Button } from 'flowbite-svelte';

	const SAMPLE_RATE = 16000;

	const transcription = new SvelteMap();

	let audioContext: AudioContext;
	onMount(async function () {
		audioContext = new AudioContext({
			sampleRate: SAMPLE_RATE
		});
		await audioContext.audioWorklet.addModule('/pcm-webaudio.js');

		const extractor = new AudioWorkletNode(audioContext, 'pcm-extractor');

		const stream = await navigator.mediaDevices.getUserMedia({ audio: { channelCount: 1 } });
		const inputAudio = audioContext.createMediaStreamSource(stream);

		let pcmBuffer = new Float32Array();
		extractor.port.onmessage = async function (evt) {
			pcmBuffer = new Float32Array([...pcmBuffer, ...new Float32Array(evt.data)]);
			console.info('Received...', 'New length:', pcmBuffer.length);
		};

		inputAudio.connect(extractor);
		extractor.connect(audioContext.destination);

		let stillProcessing = false;
		setInterval(async () => {
			if (stillProcessing) {
				return;
			}
			stillProcessing = true;

			const pcm = new Float32Array([...pcmBuffer]);
			pcmBuffer = new Float32Array();

			if (pcm.length == 0) {
				stillProcessing = false;
				return;
			}

			console.info('Sending PCM data...', 'Length:', pcm.length);

			try {
				const resp = await fetch('/broadcast/process-pcm', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/octet-stream'
					},
					body: pcm
				});

				const json = (await resp.json()).text;
				transcription.set(json.line, json.text);
			} catch (e) {
			} finally {
				stillProcessing = false;
			}
		}, 2000);
	});

	async function startRecording() {
		await audioContext.resume();
	}
	async function stopRecording() {
		await audioContext.suspend();
	}
</script>

<Button onclick={startRecording}>Start</Button>
<Button onclick={stopRecording}>Stop</Button>

<div id="transmission-text">
	{#each transcription.entries() as [line, text]}
		<span>
			<p>{text}</p>
		</span>
	{/each}
</div>
