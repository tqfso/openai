# 模型开放平台

## 目录结构

- common 共用包
- openserver 开放平台服务
- apiserver 网关服务

## 时序数据库

### 部署“TimescaleDB”

```sh
docker run -d --name timescaledb \
  -p 11900:5432 \
  -e POSTGRES_PASSWORD=123456 \
  -v /zol/postgresql/data:/var/lib/postgresql/data \
  timescale/timescaledb:latest-pg17
```

### 创建数据库

```sh
docker exec -it timescaledb psql -U postgres

CREATE DATABASE openai_db;

\c openai_db  -- 切换到新数据库

CREATE EXTENSION IF NOT EXISTS timescaledb;

```
