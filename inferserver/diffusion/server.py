import os
import sys
import time
import argparse
import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
from diffusers import DiffusionPipeline
import torch
import base64
from io import BytesIO

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

class OpenAIExtraBody(BaseModel):
    num_inference_steps: Optional[int] = 10
    denoising_end: Optional[float] = 1.0

class OpenAIImageGenerationsRequest(BaseModel):
    model: str
    prompt: str
    n: Optional[int] = 1
    size: Optional[str] = 'auto'
    output_format: Optional[Literal['png', 'jpeg', 'webp']] = 'png'
    response_format: Optional[Literal['url', 'b64_json']] = 'b64_json'
    stream: Optional[bool] = False
    extra_body: Optional[OpenAIExtraBody] = OpenAIExtraBody()

class OpenAIImageGenerationsResponse(BaseModel):
    created: int
    data: List[Dict]
    usage: Optional[Dict] = None

@app.post("/v1/images/generations")
def inference_v1(request: OpenAIImageGenerationsRequest):    
    image = generate_image(
        prompt=request.prompt,
        num_inference_steps=request.extra_body.num_inference_steps,
        denoising_end=request.extra_body.denoising_end,
    )

    buffered = BytesIO()
    image.save(buffered, format="PNG")
    img_b64 = base64.b64encode(buffered.getvalue()).decode()

    return OpenAIImageGenerationsResponse(
        created=int(time.time()),
        data=[{"b64_json": img_b64} for _ in range(request.n)],
        usage=None
    )

# 加载模型
def load_model(model_dir: str):
    base = DiffusionPipeline.from_pretrained(
        model_dir, 
        torch_dtype=torch.float16, 
        variant="fp16", 
        use_safetensors=True
    )
    base.enable_model_cpu_offload()

    refiner = DiffusionPipeline.from_pretrained(
        model_dir,
        text_encoder_2=base.text_encoder_2,
        vae=base.vae,
        torch_dtype=torch.float16,
        use_safetensors=True,
        variant="fp16",
    )
    refiner.enable_model_cpu_offload()
    return base, refiner

# 生成图片
def generate_image(prompt: str, num_inference_steps: int, denoising_end: float):
    image = base(
        prompt=prompt,
        num_inference_steps=num_inference_steps,
        denoising_end=denoising_end,
        output_type="latent",
    ).images

    image = refiner(
        prompt=prompt,
        num_inference_steps=num_inference_steps,
        denoising_start=denoising_end,
        image=image,
    ).images[0]

    return image

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--port', type=int, default=8000)
    parser.add_argument('--model_dir', type=str, default='/models/stabilityai/stable-diffusion-xl-base-1.0', help='local path or modelscope repo id')
    args = parser.parse_args()

    if not os.path.isdir(args.model_dir):  # to trigger model download if needed
        logging.info(f"Model directory {args.model_dir} does not exist locally")
        sys.exit(0)
    
    try:
        base, refiner = load_model(args.model_dir)
    except Exception:
        sys.exit(0)
    
    uvicorn.run(app, host="0.0.0.0", port=args.port)
