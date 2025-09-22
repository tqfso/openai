class BaseModel:
    def __init__(self, model_dir, dtype):
        self.model_dir = model_dir
        self.dtype = dtype

    def generate(self, model_inputs):
        raise NotImplementedError