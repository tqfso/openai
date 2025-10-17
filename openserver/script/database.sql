
-- ====================
-- 数据库初始化脚本
-- ====================

/* 系统配置表 */
DROP TABLE IF EXISTS system_configs;
CREATE TABLE system_configs (  
    key TEXT PRIMARY KEY, -- 配置项名称
    value TEXT NOT NULL, -- 配置项值
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* 拓扑域表 */
DROP TABLE IF EXISTS topo_domains;
CREATE TABLE topo_domains (
    id BIGINT PRIMARY KEY, -- 拓扑域ID
    vpc_id BIGINT NOT NULL, -- 私有网络ID
    status TEXT DEFAULT 'enabled', -- 状态: enabled, disabled
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* 推理引擎表 */
DROP TABLE IF EXISTS infer_engines;
CREATE TABLE infer_engines (
    name TEXT PRIMARY KEY, -- 推理引擎名称: vllm-openai
    framework TEXT NOT NULL, -- 推理引擎框架: Vllm, SLang, Ollama, LMDeploy, Pytorch
    image TEXT NOT NULL, -- 镜像名称
    status TEXT DEFAULT 'enabled', -- 状态: enabled, disabled
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO infer_engines(name, framework, image) VALUES 
('vllm-openai', 'Vllm', 'reasoning/vllm/vllm-openai:latest'),
('zdan-qwenvl', 'Pytorch', 'reasoning/zdan/qwenvl:latest');

/* 模型类型表 */
DROP TABLE IF EXISTS model_classes;
CREATE TABLE model_classes (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL, -- 类别名称
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO model_classes (id, name) VALUES
(1, '文本生成'),
(2, '向量嵌入'),
(3, '重排序'),
(4, '视觉理解'),
(5, '图像生成'),
(6, '视频生成'),
(7, '语音识别'),
(8, '语音合成'),
(9, '多模态模型'),
(10, '深度思考');

/* 模型提供商 */
DROP TABLE IF EXISTS model_providers;
CREATE TABLE model_providers (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO model_providers(id, name) VALUES
(1, '通义千问'),
(2, 'DeepSeek'),
(3, '月之暗面'),
(4, '智谱AI'),
(5, 'Black Forest Labs'),
(6, 'MiniMax'),
(7, 'Stability AI');

/* 模型扩展能力 */
DROP TABLE IF EXISTS model_abilities;
CREATE TABLE model_abilities (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO model_abilities(id, name) VALUES
(1, '工具调用'),
(2, '结构化输出'),
(3, 'Cache缓存'),
(4, '批量推理'),
(5, '模型体验'),
(6, '模型调优'),
(7, '联网搜索');

/* 预置模型表*/
DROP TABLE IF EXISTS platform_models;
CREATE TABLE platform_models (
    name TEXT PRIMARY KEY, -- 模型名称如 Qwen/Qwen3-Reranker-8B
    provider BIGINT NOT NULL, -- 深度求索、通义实验室等
    classes BIGINT[] NOT NULL, -- 模型类型
    abilities BIGINT[], -- 扩展能力
    max_context_length BIGINT DEFAULT 0, -- 最大上下文长度
	deploy_info JSONB, -- 部署信息：支持的推理引擎列表(推理引擎、可用加速卡、运行命令、运行参数、环境变量等)
    finetune_info JSONB, -- 微调信息: 支持的训练引擎列表(微调引擎、可用加速卡、运行命令、运行参数、环境变量等)
    description TEXT, -- 描述
    status TEXT DEFAULT 'enabled', -- 状态: enabled, disabled
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* 平台模型服务表 */
DROP TABLE IF EXISTS platform_services;
CREATE TABLE platform_services (
    id TEXT PRIMARY KEY, -- 模型服务ID，资源调度返回的服务ID
	name TEXT NOT NULL, -- 服务名称
	topo_id BIGINT NOT NULL, -- 所属拓扑域
    model_name TEXT NOT NULL, -- 模型名称
    api_service_id TEXT NOT NULL, -- API网关服务
    power BIGINT NOT NULL DEFAULT 0, -- 部署的算力
    load BIGINT NOT NULL DEFAULT 0, -- 平均负载
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* API网关服务表 */
DROP TABLE IF EXISTS api_services;
CREATE TABLE api_services (
    id TEXT PRIMARY KEY, -- 网关服务ID，资源调度返回的服务ID
    name TEXT UNIQUE NOT NULL, -- 名称
    topo_id BIGINT NOT NULL, -- 所属拓扑域
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* 用户表 */
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id TEXT PRIMARY KEY, -- 用户ID，与零极云保持一致
    nick_name TEXT , -- 昵称
    request_limit BIGINT DEFAULT 60, -- 请求数限流（次/分钟）
    token_limit BIGINT DEFAULT 1000000,  -- Token限流（Tokens/分钟）
    status TEXT DEFAULT 'enabled', -- 状态: enabled, disabled
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* 用户工作空间表 */
DROP TABLE IF EXISTS workspaces;
CREATE TABLE workspaces (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL, -- 用户ID
    name TEXT NOT NULL, -- 工作空间名称
    status TEXT DEFAULT 'enabled', -- 状态: enabled, disabled
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, name)
);

/* 工作空间调用限制 */
DROP TABLE IF EXISTS usage_limits;
CREATE TABLE usage_limits (
    workspace_id TEXT NOT NULL, -- 所属工作空间ID
    model_name TEXT NOT NULL, -- 模型服务名称
    request_limit BIGINT NOT NULL, -- 请求数限流（次/分钟）
    token_limit BIGINT NOT NULL,  -- Token限流（Tokens/分钟）
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (workspace_id, model_name)
);

/* API密钥表 */
DROP TABLE IF EXISTS api_keys;
CREATE TABLE api_keys (
    id TEXT PRIMARY KEY,  -- API密钥
    user_id TEXT NOT NULL, -- 用户ID
    workspace_id TEXT NOT NULL, -- 所属工作空间ID
    description TEXT, -- 描述
    expires_at TIMESTAMPTZ, -- 过期时间，NULL 表示永不过期
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()

);

/* 调用统计表 */
DROP TABLE IF EXISTS usage_logs;
CREATE TABLE usage_logs (
    id BIGSERIAL,
    api_key TEXT NOT NULL, -- 调用密钥
    user_id TEXT NOT NULL, -- 用户ID
    service_id BIGINT NOT NULL, -- 服务ID
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- 精确到毫秒
    status SMALLINT DEFAULT 0, -- 调用状态
    input_tokens BIGINT DEFAULT 0, -- 输入token数量
    output_tokens BIGINT DEFAULT 0, -- 输出token数量
	response_time_ms INT NOT NULL,  -- 响应耗时(毫秒)
    PRIMARY KEY (id, occurred_at)
);

SELECT create_hypertable('usage_logs', 'occurred_at');


