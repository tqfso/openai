from openai import OpenAI
from config import BASE_URL

# 测试图片分析

model_name = "Qwen/Qwen2.5-VL-7B-Instruct-AWQ"

client = OpenAI(
    base_url=BASE_URL,
    api_key="EMPTY"
)

response = client.responses.create(
    model=model_name,
    input=[
        {
            "role": "user",
            "content": [
                {
                    "type": "input_text", 
                    "text": "图片里有什么?"
                },

                {
                    "type": "input_image",
                    "image_url": "https://qianwen-res.oss-cn-beijing.aliyuncs.com/Qwen-VL/assets/demo.jpeg",
                },
            ],
        }
    ],
)

print(response.output_text)