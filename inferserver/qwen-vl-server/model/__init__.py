import os
from .qwen_vl import QwenVLModel

def get_model_name_from_path(model_dir):
    return os.path.basename(os.path.normpath(model_dir))

def create_model(model_dir, dtype):
    model_name = get_model_name_from_path(model_dir)    
    if model_name == "Qwen2.5-VL-3B-Instruct-AWQ":
        return QwenVLModel(model_dir=model_dir, dtype=dtype)
    # 可扩展支持其他模型
    raise ValueError(f"Unknown model name: {model_name}")
