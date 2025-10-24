class BaseModel:
    def __init__(self, model_path):
        self.model_path = model_path

    def generate(self, model_inputs):
        raise NotImplementedError