from modelscope import Qwen2_5_VLForConditionalGeneration, AutoProcessor
from qwen_vl_utils import process_vision_info
import torch
from model.base import BaseModel
import pdb

class QwenVLModel(BaseModel):
    def __init__(self, model_dir, dtype):
        super().__init__(model_dir, dtype)
        self.model = Qwen2_5_VLForConditionalGeneration.from_pretrained(
            model_dir, dtype=getattr(torch, dtype), device_map="auto"
        )
        self.processor = AutoProcessor.from_pretrained(model_dir)

    def generate(self, model_inputs):
        messages = model_inputs.get("messages")
        max_tokens = model_inputs.get("max_tokens", 128)
        text = self.processor.apply_chat_template(messages, tokenize=False, add_generation_prompt=True)
        image_inputs, video_inputs = process_vision_info(messages)
        inputs = self.processor(
            text=[text],
            images=image_inputs,
            videos=video_inputs,
            padding=True,
            return_tensors="pt",
        )
        inputs = inputs.to("cuda")

        # 统计输入token数（包括图片）
        input_tokens = len(self.processor.tokenizer(text)["input_ids"])

        # Inference: Generation of the output
        generated_ids = self.model.generate(**inputs, max_new_tokens=max_tokens)
        generated_ids_trimmed = [
            out_ids[len(in_ids) :] for in_ids, out_ids in zip(inputs.input_ids, generated_ids)
        ]
        output_text = self.processor.batch_decode(
            generated_ids_trimmed, skip_special_tokens=True, clean_up_tokenization_spaces=False
        )

        # 统计输出token数
        output_tokens = len(self.processor.tokenizer(output_text[0])["input_ids"]) if output_text else 0

        return {
            "output": output_text[0] if output_text else "",
            "usage": {
                "input_tokens": input_tokens,
                "output_tokens": output_tokens,
                "total_tokens": input_tokens + output_tokens
            }
        }