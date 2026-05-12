<script lang="ts">
	import { onMount } from 'svelte';

	import { Timeline, TimelineItem, Button, Modal } from 'flowbite-svelte';
	import { APIProvider, GoogleMap, Marker } from 'svelte-google-maps-api';

	import { uuidV4AsString } from '$lib/uuid';
	import { env } from '$env/dynamic/public';

	const SAMPLE_RATE = 16000;
	const SAMPLE_UPLOAD_WAIT = 1000; // In milliseconds.

	const CURRENT_SESSION_ID = uuidV4AsString();

	const data: any = $props();

	const apiKey = env.PUBLIC_GOOGLE_MAPS_API_KEY;
	const mapOptions: { zoom: number; center: { lat: number; lng: number } | null } = $state({
		zoom: 14,
		center: null
	});

	onMount(() => {
		navigator.geolocation.getCurrentPosition(
			(position) => {
				mapOptions.center = {
					lat: position.coords.latitude,
					lng: position.coords.longitude
				};
			},
			(error) => {
				console.error('Error getting location:', error);
			}
		);
	});

	let showMap = $state(false);

	let lines = $state([]);

	let audioContext: AudioContext;
	onMount(async function () {
		// TODO: Handle deny of the mic.
		audioContext = new AudioContext({
			sampleRate: SAMPLE_RATE
		});
		await audioContext.audioWorklet.addModule('/pcm-webaudio.js');

		const extractor = new AudioWorkletNode(audioContext, 'pcm-extractor');

		const stream = await navigator.mediaDevices.getUserMedia({ audio: { channelCount: 1 } });
		const inputAudio = audioContext.createMediaStreamSource(stream);

		// TODO: Consider working with `ArrayBuffer` directly.
		const pcmQueue: Float32Array[] = [];

		extractor.port.onmessage = async function (evt) {
			pcmQueue.push(new Float32Array(evt.data));
		};

		inputAudio.connect(extractor);
		extractor.connect(audioContext.destination);

		const sendSamples = async () => {
			const samplesToSend = pcmQueue.length;
			if (samplesToSend == 0) {
				setTimeout(sendSamples, SAMPLE_UPLOAD_WAIT);
				return;
			}

			let bufferSize = 0;
			for (let i = 0; i < samplesToSend; i++) {
				bufferSize += pcmQueue[i].length;
			}

			const buffer = new Float32Array(bufferSize);
			let bufferOffset = 0;
			for (let i = 0; i < samplesToSend; i++) {
				const pcm = pcmQueue.shift()!;
				buffer.set(pcm, bufferOffset);
				bufferOffset += pcm.length;
			}

			console.assert(buffer.length > 0, 'The buffer is empty!');

			try {
				const req = await fetch('/broadcast/process-pcm', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/octet-stream',
						'X-Current-Session-Id': CURRENT_SESSION_ID
					},
					body: buffer
				});

				lines = await req.json();
			} catch (e) {
				console.log('The server rejected the PCM data!', e);
			}

			setTimeout(sendSamples, SAMPLE_UPLOAD_WAIT);
		};

		sendSamples();
	});

	async function startRecording() {
		await audioContext.resume();
	}
	async function stopRecording() {
		await audioContext.suspend();
	}
</script>

<div class="flex h-full flex-col">
	<!--
	<Timeline class="m-8 grow overflow-y-scroll">
		TODO: Fix this.
		{#each lines as line}
			<TimelineItem title={line} date=""><p></p></TimelineItem>
		{/each}
	</Timeline>
		-->

	{#if mapOptions.center}
		<Modal
			bind:open={showMap}
			onclose={() => {
				showMap = false;
			}}
		>
			<div style="width: 600px; height: 400px;">
				<APIProvider {apiKey}>
					<GoogleMap options={mapOptions} mapContainerStyle="width: 100%; height: 100%;" />
				</APIProvider>
			</div>
		</Modal>
	{/if}

	<div class="m-4 flex justify-center gap-4">
		<Button onclick={startRecording}>Transmitir</Button>
		<Button onclick={stopRecording}>Detener</Button>
		{#if mapOptions.center}
			<Button
				onclick={() => {
					showMap = true;
				}}>Mostrar mapa</Button
			>
		{/if}
	</div>
</div>
