<script lang="ts">
import { SvelteMap } from "svelte/reactivity";
import { onMount } from "svelte";
import { Button } from "flowbite-svelte";

const SAMPLE_RATE = 16000;

const transcription = new SvelteMap();

let audioContext: AudioContext;
onMount(async function () {
	audioContext = new AudioContext({
		sampleRate: SAMPLE_RATE,
	});
	await audioContext.audioWorklet.addModule("/pcm-webaudio.js");

	const extractor = new AudioWorkletNode(audioContext, "pcm-extractor");

	const stream = await navigator.mediaDevices.getUserMedia({
		audio: { channelCount: 1 },
	});
	const inputAudio = audioContext.createMediaStreamSource(stream);

	let pcmBuffer = new Float32Array();
	let processing = false;

	extractor.port.onmessage = async function (evt) {
		const audioData = new Float32Array(evt.data);

		// Efficient accumulation using set() instead of spread
		const newBuffer = new Float32Array(pcmBuffer.length + audioData.length);
		if (pcmBuffer.length > 0) {
			newBuffer.set(pcmBuffer, 0);
		}
		newBuffer.set(audioData, pcmBuffer.length);
		pcmBuffer = newBuffer;

		if (pcmBuffer.length >= 16000 * 2 && !processing) {
			processing = true;
			const toSend = pcmBuffer.slice(0, 16000 * 2);

			// Keep overflow efficiently
			if (pcmBuffer.length > 16000 * 2) {
				pcmBuffer = pcmBuffer.slice(16000 * 2);
			} else {
				pcmBuffer = new Float32Array();
			}

			console.info("Sending PCM data...", "Length:", toSend.length);

			try {
				const resp = await fetch("/broadcast/process-pcm", {
					method: "POST",
					headers: {
						"Content-Type": "application/octet-stream",
					},
					body: toSend,
				});

				if (resp.ok) {
					const json = (await resp.json()).text;
					if (json.text && json.text.trim() !== "") {
						transcription.set(json.line, json.text);
					}
				}
			} catch (e) {
				console.error("Transcription error:", e);
			} finally {
				processing = false;
			}
		}
	};

	inputAudio.connect(extractor);
	extractor.connect(audioContext.destination);
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
