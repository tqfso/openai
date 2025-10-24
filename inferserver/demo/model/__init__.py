import os
from .qwen_vl import QwenVLModel

def get_model_name_from_path(model_dir):
    return os.path.basename(os.path.normpath(model_dir))

def create_model(model_dir):
    return QwenVLModel(model_dir=model_dir)
