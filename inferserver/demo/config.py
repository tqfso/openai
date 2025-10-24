import os

class Config:
    MODEL_DIR = os.getenv("MODEL_DIR")

config = Config()