import os
import argparse
from fastapi import FastAPI
from model import create_model
from config import config
from api import chat_completions_router

# 启用碎片优化
os.environ["PYTORCH_CUDA_ALLOC_CONF"] = "expandable_segments:True"

# 读取命令参数
parser = argparse.ArgumentParser(description="服务启动参数")
parser.add_argument("--model", type=str, help="模型名称", default="Qwen2.5-VL-3B-Instruct-AWQ")
args = parser.parse_args()

# 初始化模型
config.model_name = args.model
model = create_model(args.model)

# 创建 FastAPI 应用
app = FastAPI()
app.state.model = model
app.include_router(chat_completions_router)

@app.get("/")
def health():
    return {"status": "ok"}