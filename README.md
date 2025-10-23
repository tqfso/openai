# 模型开放平台

## 目录结构

- common 共用包
- openserver 开放平台服务
- apiserver API网关服务
- inferserver 推理服务
- trainserver 训练服务

## 时序数据库

### 部署

```sh
docker run -d --name timescaledb \
        -p 11900:5432 \
        -e POSTGRES_PASSWORD=123456 \
        -v /zol/postgresql/data:/var/lib/postgresql/data \
        timescale/timescaledb:latest-pg17
```

### 命令

```sh

psql -U postgres -- 以管理员账号进入命令交互

[pg] CREATE DATABASE openai_db; -- 创建数据库

[pg] \c openai_db  -- 切换到新数据库

[pg] CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

```

## 模型

### Qwen/Qwen3-1.7B

- download from modescope

```sh
export MODELSCOPE_CACHE="/zol"
modelscope download --model Qwen/Qwen3-1.7B
cd /zol/models
rm Qwen3-1.7B
mv Qwen3-1___7B Qwen3-1.7B

```

- convert to Harmony format

```sh
export GIT_SSL_NO_VERIFY=0
git clone --recursive https://github.com/mlc-ai/mlc-llm.git
```

- deploy using Tesla T4
  
```sh
docker run --rm --gpus all --shm-size=1g \
	-p 8000:8000 \
	-v /zol/models:/vllm-workspace \
	-e VLLM_USE_FLASHINFER_SAMPLER=0 \
	--name qwen3-1.7b \
	vllm/vllm-openai:0.11.0 \
	--model Qwen/Qwen3-1.7B \
	--enable-auto-tool-choice \
	--tool-call-parser hermes \
	--dtype float32 \
	--max-model-len 2800 \
	--gpu-memory-utilization 0.80


```
### Qwen/Qwen2.5-VL-3B-Instruct-AWQ

- deploy sing Tesla T4

```sh
docker run --rm --gpus all --shm-size=1g \
	-p 8000:8000 \
	-v /zol/models:/vllm-workspace \
	-e VLLM_USE_FLASHINFER_SAMPLER=0 \
	--name qwen25-vl \
	vllm/vllm-openai:0.11.0 \
	--model Qwen/Qwen2.5-VL-3B-Instruct-AWQ \
	--mm-processor-kwargs '{"max_pixels": 262144}' \
	--max-model-len 2048 \
	--gpu-memory-utilization 0.80
```