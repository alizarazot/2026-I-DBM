import struct
import json
import time
import wave
import os
import signal
import atexit
import sys
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from collections import defaultdict

from moonshine_voice import (
    Transcriber,
    get_model_for_language,
)

SAMPLE_RATE = 16000
DEBUG = True
DEBUG_DIR = Path("debug_audio")
server_instance = None

if DEBUG and not DEBUG_DIR.exists():
    DEBUG_DIR.mkdir(parents=True, exist_ok=True)


class Server(BaseHTTPRequestHandler):
    sessions = {}
    audio_counter = 0
    session_audio_data = defaultdict(list)  # session_id -> list of audio chunks
    session_start_time = {}  # session_id -> start timestamp
    
    def do_DELETE(self):
        """Clean up a session when client disconnects."""
        session_id = self.headers.get("X-Session-Id")
        if not session_id:
            self.send_response(401)
            self.end_headers()
            return
        
        if session_id in self.sessions:
            # Save final merged audio before cleanup
            if DEBUG:
                self._save_final_merged_audio(session_id)
            
            # Clean up transcriber
            transcriber = self.sessions[session_id]
            transcriber.stop()
            del self.sessions[session_id]
            
            # Clean up debug data
            if session_id in self.session_audio_data:
                del self.session_audio_data[session_id]
            if session_id in self.session_start_time:
                del self.session_start_time[session_id]
            
            print(f"DEBUG: Cleaned up session {session_id}")
        
        self.send_response(200)
        self.end_headers()

    def do_POST(self):
        session_id = self.headers.get("X-Session-Id")
        if not session_id:
            self.send_response(401)
            self.end_headers()
            return

        content_length = int(self.headers.get("Content-Length", 0))
        audio_raw = self.rfile.read(content_length)
        
        # Debug: Save incoming audio
        if DEBUG:
            timestamp = int(time.time() * 1000)
            
            # Save individual chunk
            chunk_filename = DEBUG_DIR / f"audio_{session_id}_{timestamp}_{Server.audio_counter}.wav"
            Server.audio_counter += 1
            
            with wave.open(str(chunk_filename), 'wb') as wav_file:
                wav_file.setnchannels(1)
                wav_file.setsampwidth(4)  # 32-bit float = 4 bytes
                wav_file.setframerate(SAMPLE_RATE)
                wav_file.writeframes(audio_raw)
            
            print(f"DEBUG: Saved chunk {content_length} bytes to {chunk_filename}")
            print(f"DEBUG: Sample count: {content_length // 4}")
            
            # Accumulate for merged file
            Server.session_audio_data[session_id].append(audio_raw)
            
            # Initialize session start time if not set
            if session_id not in Server.session_start_time:
                Server.session_start_time[session_id] = timestamp

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

        transcription = transcriber.update_transcription()
        if transcription.lines:
            line = transcription.lines[-1]
            print(line.start_time, line.text)
            response_text = line.text
        else:
            response_text = ""
        
        # Debug: Save merged audio periodically
        if DEBUG:
            self._save_merged_audio(session_id, timestamp)
        
        self.send_response(200)
        self.send_header("Content-Type", "text/json")
        self.end_headers()
        
        self.wfile.write(
            bytes(json.dumps({"line": 0 if not transcription.lines else line.start_time, "text": response_text}), "utf-8")
        )
    
    def _save_merged_audio(self, session_id, current_timestamp):
        """Save merged audio for a session every 5 seconds or when session ends."""
        if session_id not in Server.session_audio_data:
            return
        
        # Save merged file every 5 seconds or if we have more than 10 chunks
        session_start = Server.session_start_time[session_id]
        time_elapsed = (current_timestamp - session_start) / 1000  # seconds
        
        if time_elapsed >= 5 or len(Server.session_audio_data[session_id]) >= 10:
            merged_filename = DEBUG_DIR / f"merged_{session_id}_{session_start}_{current_timestamp}.wav"
            
            # Concatenate all audio chunks
            merged_data = b''.join(Server.session_audio_data[session_id])
            
            with wave.open(str(merged_filename), 'wb') as wav_file:
                wav_file.setnchannels(1)
                wav_file.setsampwidth(4)  # 32-bit float = 4 bytes
                wav_file.setframerate(SAMPLE_RATE)
                wav_file.writeframes(merged_data)
            
            print(f"DEBUG: Saved merged {len(merged_data)} bytes ({len(merged_data)//4} samples) to {merged_filename}")
            print(f"DEBUG: Created from {len(Server.session_audio_data[session_id])} chunks")
            
            # Clear accumulated data but keep tracking
            Server.session_audio_data[session_id] = []
    
    def _save_final_merged_audio(self, session_id):
        """Save final merged audio when session ends."""
        if session_id not in Server.session_audio_data:
            return
        
        if Server.session_audio_data[session_id]:
            current_timestamp = int(time.time() * 1000)
            session_start = Server.session_start_time.get(session_id, current_timestamp)
            merged_filename = DEBUG_DIR / f"final_merged_{session_id}_{session_start}_{current_timestamp}.wav"
            
            # Concatenate all audio chunks
            merged_data = b''.join(Server.session_audio_data[session_id])
            
            with wave.open(str(merged_filename), 'wb') as wav_file:
                wav_file.setnchannels(1)
                wav_file.setsampwidth(4)  # 32-bit float = 4 bytes
                wav_file.setframerate(SAMPLE_RATE)
                wav_file.writeframes(merged_data)
            
            print(f"DEBUG: Saved FINAL merged {len(merged_data)} bytes ({len(merged_data)//4} samples) to {merged_filename}")
            print(f"DEBUG: Created from {len(Server.session_audio_data[session_id])} chunks")
    
    def log_request(self, code='-', size='-'):
        """Override to reduce log noise."""
        if DEBUG:
            super().log_request(code, size)
        # Otherwise silent


def cleanup():
    """Clean up all sessions on exit."""
    if DEBUG:
        print("\nDEBUG: Server shutting down, saving final merged files...")
        for session_id in Server.sessions.copy():
            Server.sessions[session_id].stop()
        
        # Save any remaining audio data
        for session_id in list(Server.session_audio_data.keys()):
            if Server.session_audio_data[session_id]:
                current_timestamp = int(time.time() * 1000)
                session_start = Server.session_start_time.get(session_id, current_timestamp)
                merged_filename = DEBUG_DIR / f"exit_merged_{session_id}_{session_start}_{current_timestamp}.wav"
                
                merged_data = b''.join(Server.session_audio_data[session_id])
                
                with wave.open(str(merged_filename), 'wb') as wav_file:
                    wav_file.setnchannels(1)
                    wav_file.setsampwidth(4)
                    wav_file.setframerate(SAMPLE_RATE)
                    wav_file.writeframes(merged_data)
                
                print(f"DEBUG: Saved exit merged {len(merged_data)} bytes to {merged_filename}")


def run(handler_class=Server):
    atexit.register(cleanup)
    signal.signal(signal.SIGINT, lambda sig, frame: exit(0))
    signal.signal(signal.SIGTERM, lambda sig, frame: exit(0))
    
    server_address = "localhost", 8000
    httpd = ThreadingHTTPServer(server_address, handler_class)
    global server_instance
    server_instance = httpd
    
    print(f"DEBUG: Starting server on {server_address}")
    print(f"DEBUG: Saving audio files to {DEBUG_DIR.absolute()}")
    print("DEBUG: Press Ctrl+C to stop server and save final merged files")
    httpd.serve_forever()


if __name__ == "__main__":
    run()
