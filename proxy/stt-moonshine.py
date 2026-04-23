from typing import Annotated, List
from dataclasses import dataclass

import struct

from fastapi import FastAPI, Header, Body
from moonshine_voice import (
    Transcriber,
    TranscriptEventListener,
    get_model_for_language,
)

SAMPLE_RATE = 16000


class TTSEventListener(TranscriptEventListener):
    def __init__(self, lines_list_reference: List[str], *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.lines = lines_list_reference

    def on_line_completed(self, event):
        self.lines.append(event.line.text)
        print(f"Line completed: {event.line.text}")


@dataclass
class TranscriberWithEvent:
    transcriber: Transcriber
    lines: List[str]


transcribers = {}

model_path, model_arch = get_model_for_language("es")


app = FastAPI()


# TODO: Concurrent requests to the same session can break this.
@app.post("/")
async def post_root(
    x_session_id: Annotated[str, Header()],
    x_current_session_id: Annotated[str, Header()],
    data: bytes = Body(),
):
    key = f"{x_session_id}-{x_current_session_id}"
    if key not in transcribers:
        transcribers[key] = TranscriberWithEvent(
            Transcriber(model_path=model_path, model_arch=model_arch), list()
        )
        transcriber = transcribers[key]

        transcriber.transcriber.add_listener(TTSEventListener(transcribers[key].lines))

        transcriber.transcriber.start()

    transcriber = transcribers[key]

    data_size = len(data) // 4  # Is a JS Float32Array.
    samples = struct.unpack(f"<{data_size}f", data)

    transcriber.transcriber.add_audio(samples, SAMPLE_RATE)

    return transcriber.lines
