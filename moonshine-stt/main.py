import struct
import json
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

from moonshine_voice import (
    Transcriber,
    get_model_for_language,
)

SAMPLE_RATE = 16000


class Server(BaseHTTPRequestHandler):
    sessions = {}

    def do_POST(self):
        session_id = self.headers.get("X-Session-Id")
        if not session_id:
            self.send_response(401)
            self.end_headers()
            return

        content_length = int(self.headers.get("Content-Length", 0))
        audio_raw = self.rfile.read(content_length)

        num_floats = len(audio_raw) // 4  # Browsers save Float32Array as 4 bytes.
        audio_data = list(struct.unpack(f"<{num_floats}f", audio_raw))

        if session_id not in self.sessions:
            model_path, model_arch = get_model_for_language("es", None)

            transcriber = Transcriber(model_path=model_path, model_arch=model_arch)

            transcriber.remove_all_listeners()
            transcriber.start()

            self.sessions[session_id] = transcriber

        transcriber = self.sessions[session_id]
        transcriber.add_audio(audio_data, SAMPLE_RATE)

        line = transcriber.update_transcription().lines[-1]
        print(line.start_time, line.text)

        self.send_response(200)

        self.send_header("Content-Type", "text/json")
        self.end_headers()

        self.wfile.write(
            bytes(json.dumps({"line": line.start_time, "text": line.text}), "utf-8")
        )


def run(handler_class=Server):
    server_address = "localhost", 8000
    ThreadingHTTPServer(server_address, handler_class).serve_forever()


if __name__ == "__main__":
    run()
