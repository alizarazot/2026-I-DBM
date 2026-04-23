class PcmExtractor extends AudioWorkletProcessor {
	constructor() {
		super();
		this.buffer = new Float32Array(64000); // 4 seconds capacity for overflow
		this.writePos = 0;
		this.readPos = 0;
		this.available = 0;
	}

	process(inputs) {
		const channel = inputs[0][0];
		if (!channel) return true;

		const inputLength = channel.length;

		// Copy all input into circular buffer (overwrites old data if needed)
		for (let i = 0; i < inputLength; i++) {
			this.buffer[this.writePos] = channel[i];
			this.writePos = (this.writePos + 1) % this.buffer.length;

			// If buffer is full, overwrite from beginning
			if (this.available >= this.buffer.length) {
				this.readPos = (this.readPos + 1) % this.buffer.length;
			} else {
				this.available++;
			}
		}

		// Send if we have at least 2 seconds (32k samples)
		if (this.available >= 32000) {
			const toSend = new Float32Array(32000);

			// Read 32000 samples from circular buffer
			let tempReadPos = this.readPos;
			for (let i = 0; i < 32000; i++) {
				toSend[i] = this.buffer[tempReadPos];
				tempReadPos = (tempReadPos + 1) % this.buffer.length;
			}

			this.port.postMessage(toSend.buffer, [toSend.buffer]);

			// Update positions
			this.readPos = (this.readPos + 32000) % this.buffer.length;
			this.available -= 32000;
		}

		return true;
	}
}

registerProcessor("pcm-extractor", PcmExtractor);
