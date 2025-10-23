from openai import OpenAI
from config import BASE_URL

model_name = "Qwen/Qwen3-1.7B"

client = OpenAI(
    base_url=BASE_URL,
    api_key="EMPTY"
)

messages = [
    {"role": "system", "content": "你是一名气象助手"},
    {"role": "user", "content": "今天的天气怎么样，共有哪几种描述?"}
]

response = client.chat.completions.create(
    model=model_name,
    messages=messages,
    max_tokens=1024,
    temperature=0.2,
    top_p=0.95,
)

print(response.choices[0].message.content)