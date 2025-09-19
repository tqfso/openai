FROM pytorch/pytorch:2.3.1-cuda12.1-cudnn8-devel

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

RUN apt-get update && apt-get install -y git && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN pip install --upgrade pip
# 安装 transformers, accelerate, fastapi, uvicorn, pillow 等
RUN pip install "transformers>=4.37" accelerate torch torchvision torchaudio --extra-index-url https://download.pytorch.org/whl/cu121
RUN pip install fastapi uvicorn pillow requests

WORKDIR /app

COPY app.py /app/app.py