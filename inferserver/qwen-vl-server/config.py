import os

class Config:
    MODEL_DIR = os.getenv("MODEL_DIR", "/work/models/Qwen/Qwen2.5-VL-3B-Instruct-AWQ")
    DTYPE = os.getenv("DTYPE", "float16")
    HOST = os.getenv("HOST", "0.0.0.0")
    PORT = int(os.getenv("PORT", "8000"))

config = Config()