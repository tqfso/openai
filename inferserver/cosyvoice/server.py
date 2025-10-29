import os
import sys
import argparse
import logging
logging.getLogger('matplotlib').setLevel(logging.WARNING)
from fastapi import FastAPI
from fastapi.responses import StreamingResponse
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import numpy as np
from cosyvoice.cli.cosyvoice import CosyVoice2
from cosyvoice.utils.file_utils import load_wav

from pydantic import BaseModel
from typing import List, Union, Literal, Optional, Dict

app = FastAPI()
# set cross region allowance
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"])

class OpenSpeechRequest(BaseModel):
    model: str
    input: str
    speed: Optional[float] = 1.0
    instructions: Optional[str] = None
    stream: Optional[bool] = False

import struct
def generate_data(model_output, sample_rate=16000):
    # 发送 WAV header（data_size 设为最大，客户端可兼容）
    def wav_header(num_channels, sample_width, sample_rate):
        data_size = 0x7FFFFFFF
        header = struct.pack(
            '<4sI4s4sIHHIIHH4sI',
            b'RIFF',
            36 + data_size,
            b'WAVE',
            b'fmt ',
            16,
            1,
            num_channels,
            sample_rate,
            sample_rate * num_channels * sample_width,
            num_channels * sample_width,
            sample_width * 8,
            b'data',
            data_size
        )
        return header
    yield wav_header(1, 2, sample_rate)
    # 实时发送PCM数据
    for tts_speech in model_output:
        yield (tts_speech['tts_speech'].numpy() * (2 ** 15)).astype(np.int16).tobytes()

@app.post("/v1/audio/speech")
def inference_v1(request: OpenSpeechRequest):    
    result = cosyvoice.inference_instruct2(request.input, request.instructions, prompt_speech_16k, speed=request.speed, stream=request.stream)
    return StreamingResponse(generate_data(result, cosyvoice.sample_rate), media_type="audio/wav")

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--port',
                        type=int,
                        default=8000)
    parser.add_argument('--model_dir',
                        type=str,
                        default='/workspace/hub/iic/CosyVoice2-0.5B',
                        help='local path or modelscope repo id')
    args = parser.parse_args()

    if not os.path.isdir(args.model_dir):  # to trigger model download if needed
        logging.info(f"Model directory {args.model_dir} does not exist locally")
        sys.exit(0)
    
    ROOT_DIR = os.path.dirname(os.path.abspath(__file__))
    sys.path.append(f'{ROOT_DIR}/third_party/Matcha-TTS')    
    prompt_speech_text = "希望你以后能够做的比我还好呦。"
    prompt_speech_16k = load_wav(f'{ROOT_DIR}/asset/zero_shot_prompt.wav', 16000)

    try:
        cosyvoice = CosyVoice2(args.model_dir)
    except Exception:
        sys.exit(0)
    
    uvicorn.run(app, host="0.0.0.0", port=args.port)
