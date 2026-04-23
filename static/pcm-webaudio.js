const SAMPLE_RATE = 16000;

// Storing a full sample rate is equal to storing whole second of audio.
const BUFFER_SECONDS = 1;
const BUFFER_SIZE = SAMPLE_RATE * BUFFER_SECONDS;

class PcmExtractor extends AudioWorkletProcessor {
	constructor() {
		super();

		this.pcmBuffer = new Float32Array(BUFFER_SIZE);
		this.pcmBufferIdx = 0;
	}

	process(inputs) {
		if (inputs.length < 1) {
			throw new Error('I need an input to process');
		}

		// TODO: Check for the other input (if any).
		if (inputs[0].length < 1) {
			// TODO: Some times the browser doesn't expose the channels, although there are ones. Find a better way to do this.
			console.error("The first input doesn't have a channel");
			return true;
		}

		const pcm = inputs[0][0];

		const remainingSpace = this.pcmBuffer.length - this.pcmBufferIdx;

		if (pcm.length <= remainingSpace) {
			this.pcmBuffer.set(pcm, this.pcmBufferIdx);
			this.pcmBufferIdx += pcm.length;
			return true;
		}

		// Fill rest of the array and send it.

		this.pcmBuffer.set(pcm.subarray(0, remainingSpace), this.pcmBufferIdx);
		this.port.postMessage(this.pcmBuffer.buffer, [this.pcmBuffer.buffer]);

		this.pcmBuffer = new Float32Array(BUFFER_SIZE); // TODO: Investigate if this allocation could cause sample drops.

		// TODO: `pcm` could still be larger than the buffer if the spec changes or the buffer is very small (<128 atm).
		this.pcmBufferIdx = pcm.length - remainingSpace;
		this.pcmBuffer.set(pcm.subarray(remainingSpace, pcm.length), 0);

		return true;
	}
}

registerProcessor('pcm-extractor', PcmExtractor);
