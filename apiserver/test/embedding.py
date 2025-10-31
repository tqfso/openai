from openai import OpenAI
from config import BASE_URL
from config import API_KEY

import numpy as np

model_name = "Qwen/Qwen3-Embedding-4B"

client = OpenAI(
    base_url=BASE_URL,
    api_key=API_KEY
)

# ===== 输入三个句子 =====
sentences = [
    "我今天心情很好",   # 目标句
    "今天的天气真不错",
    "我明天要去上班"
]

# ===== 一次性获取全部 embedding =====
response = client.embeddings.create(
    model="Qwen/Qwen3-Embedding-4B",
    input=sentences
)

# 转换为 numpy 向量
embeddings = [np.array(d.embedding) for d in response.data]

# ===== 计算余弦相似度 =====
def cosine_similarity(a, b):
    return np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b))

query_emb = embeddings[0]
candidate_embs = embeddings[1:]
candidates = sentences[1:]

similarities = [cosine_similarity(query_emb, e) for e in candidate_embs]

# ===== 输出结果 =====
for text, score in zip(candidates, similarities):
    print(f"句子：{text}\n相似度：{score:.4f}\n")

best_match = candidates[np.argmax(similarities)]
print(f"✅ 和『{sentences[0]}』最相似的句子是：『{best_match}』")