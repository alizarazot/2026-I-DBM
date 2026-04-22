
let pcmBuffer = new Float32Array();

class PcmExtractor extends AudioWorkletProcessor {
    process(inputs) {
        const channel = inputs[0][0];

        try {
            pcmBuffer = new Float32Array([...pcmBuffer, ...channel])
        } catch (e) {
            return true
        }

        if (pcmBuffer.length >= 16000 * 2 ) {
            console.info("Sending...", "Length:", pcmBuffer.length);
            this.port.postMessage(pcmBuffer.buffer, [pcmBuffer.buffer])
            pcmBuffer = new Float32Array();
        }

        return true
    }
}

registerProcessor("pcm-extractor", PcmExtractor)
