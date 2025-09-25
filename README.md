# 模型开放平台

## 目录结构

- common 共用包
- openserver 开放平台服务
- apiserver API网关服务
- reasonserver 推理服务

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

### 下载

- downloade Qwen/Qwen3-1.7B from modescope

```sh

export MODELSCOPE_CACHE="/zol"

modelscope download --model Qwen/Qwen3-1.7B

cd /zol/models
rm Qwen3-1.7B
mv Qwen3-1___7B Qwen3-1.7B

```

### 部署

- deploy Qwen/Qwen3-1.7B using Tesla T4
  
```sh
docker run --gpus all --shm-size=1g -p 8000:8000 -v /zol/models:/vllm-workspace --name qwen3-1.7b vllm/vllm-openai:latest --model Qwen/Qwen3-1.7B --enable-auto-tool-choice --tool-call-parser hermes --dtype float32 --trust-remote-code --max-model-len 2800 --gpu-memory-utilization 0.95

```
