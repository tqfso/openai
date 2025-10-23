from openai import OpenAI
from config import BASE_URL

# 测试图片分析

model_name = "Qwen/Qwen2.5-VL-3B-Instruct-AWQ"

client = OpenAI(
    base_url=BASE_URL,
    api_key="EMPTY"
)

response = client.chat.completions.create(
    model=model_name,
    messages=[
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "这张图片里有没有女孩，简单回答有或没有"
                },
                {
                    "type": "image_url",
                    "image_url": {
                        "url": "https://goal-xuxiayu.oss-cn-shanghai.aliyuncs.com/zdan/demo.jpeg"
                    }
                }
            ]
        }
        
    ]
)

print(response.choices[0].message.content)