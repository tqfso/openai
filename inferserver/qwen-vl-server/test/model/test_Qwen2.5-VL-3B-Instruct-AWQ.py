from modelscope import Qwen2_5_VLForConditionalGeneration, AutoProcessor
from qwen_vl_utils import process_vision_info
from PIL import Image
from io import BytesIO
import torch
import os
import requests

# 启用碎片优化
os.environ["PYTORCH_CUDA_ALLOC_CONF"] = "expandable_segments:True"

# 加载模型
model = Qwen2_5_VLForConditionalGeneration.from_pretrained(
    "/work/models/Qwen/Qwen2.5-VL-3B-Instruct-AWQ",
    dtype=torch.float16,
    device_map="auto",
    # attn_implementation="flash_attention_2",  # 如果支持可以取消注释
)

# 加载处理器，限制图像分辨率范围可节省显存
min_pixels = 256 * 28 * 28
max_pixels = 512 * 28 * 28
processor = AutoProcessor.from_pretrained(
    "/work/models/Qwen/Qwen2.5-VL-3B-Instruct-AWQ",
    min_pixels=min_pixels,
    max_pixels=max_pixels
)

# 加载图片（可替换为本地图像）
img_url = "https://qianwen-res.oss-cn-beijing.aliyuncs.com/Qwen-VL/assets/demo.jpeg"
image = Image.open(BytesIO(requests.get(img_url).content)).convert("RGB")
image = image.resize((448, 448))  # 推荐尺寸，防OOM

messages = [
    {
        "role": "user",
        "content": [
            {
                "type": "image",
                "image": image,
            },
            {"type": "text", "text": "描述一下这张图片"},
        ],
    }
]

# Preparation for inference
text = processor.apply_chat_template(
    messages, tokenize=False, add_generation_prompt=True
)
image_inputs, video_inputs = process_vision_info(messages)
inputs = processor(
    text=[text],
    images=image_inputs,
    videos=video_inputs,
    padding=True,
    return_tensors="pt",
)
inputs = inputs.to("cuda")

# Inference: Generation of the output
generated_ids = model.generate(**inputs, max_new_tokens=128)
generated_ids_trimmed = [
    out_ids[len(in_ids) :] for in_ids, out_ids in zip(inputs.input_ids, generated_ids)
]
output_text = processor.batch_decode(
    generated_ids_trimmed, skip_special_tokens=True, clean_up_tokenization_spaces=False
)
print(output_text)