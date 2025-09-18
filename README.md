# 模型开放平台

## 目录结构

- common 共用包
- openserver 开放平台服务
- apiserver 网关服务

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

```

## 模型

### 下载

```sh

export MODELSCOPE_CACHE="/zol"

modelscope download --model Qwen/Qwen3-1.7B

cd /zol/models
rm Qwen3-1.7B
mv Qwen3-1___7B Qwen3-1.7B

```

### 部署

```sh
docker run --gpus all \
  --shm-size=1g \
  -p 8000:8000 \
  -v /zol/models:/models \
  -d --name qwen317 \
  vllm/vllm-openai:latest \
  --model /models/Qwen/Qwen3-1.7B \
  --enable-reasoning --reasoning-parser deepseek_r1

```
