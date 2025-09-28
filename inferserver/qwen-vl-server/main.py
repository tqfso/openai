import os
from fastapi import FastAPI
from model import create_model
from config import config
from api.openai import router as openai_router
from utils.logger import setup_logger

# 启用碎片优化
os.environ["PYTORCH_CUDA_ALLOC_CONF"] = "expandable_segments:True"

# 配置日志
# setup_logger()

# 初始化模型
model = create_model(model_dir=config.MODEL_DIR, dtype=config.DTYPE)

# 创建 FastAPI 应用
app = FastAPI()
app.state.model = model
app.include_router(openai_router, prefix="/v1")

@app.get("/")
def health():
    return {"status": "ok"}