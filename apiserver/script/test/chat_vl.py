from openai import OpenAI
from config import BASE_URL

# 测试图片分析

model_name = "Qwen/Qwen2.5-VL-3B-Instruct-AWQ"

client = OpenAI(
    base_url=BASE_URL,
    api_key="EMPTY"
)

messages=[
    {
        "role": "user",
        "content": [
            {
                "type": "text",
                "text": "这是啥？"
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

response = client.chat.completions.create(
    model=model_name,
    messages=messages
)

print(response.choices[0].message.content)