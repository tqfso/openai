
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
    framework TEXT NOT NULL, -- 推理引擎框架: vllm, slang, ollama, LMDeploy, pt
    image TEXT NOT NULL, -- 镜像名称
    status TEXT DEFAULT 'enabled', -- 状态: enabled, disabled
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* 预置模型表*/
DROP TABLE IF EXISTS platform_models;
CREATE TABLE platform_models (
    name TEXT PRIMARY KEY, -- 模型名称如 Qwen/Qwen3-Reranker-8B
    provider TEXT, -- 深度求索、通义实验室等
    classes TEXT[] NOT NULL, -- 模型分类['文本生成', '图片生成', '语音识别']
    extended_ability TEXT[], -- 扩展能力如: ['tools', 'thinking', 'struct', 'batch', 'prompt caching']
    max_context_length BIGINT NOT NULL, -- 最大上下文长度
	deploy_info JSONB, -- 部署信息：支持的推理引擎列表(推理引擎、可用加速卡、运行命令、运行参数、环境变量等)
    finetune_info JSONB, -- 微调信息: 支持的训练引擎列表(微调引擎、可用加速卡、运行命令、运行参数、环境变量等)
    status TEXT DEFAULT 'enabled', -- 状态: enabled, disabled
    description TEXT, -- 描述
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* 模型服务表 */
DROP TABLE IF EXISTS model_services;
CREATE TABLE model_services (
    id TEXT PRIMARY KEY, -- 模型服务ID，资源调度返回的服务ID
	name TEXT NOT NULL, -- 服务名称
	topo_id BIGINT NOT NULL, -- 所属拓扑域
    model_name TEXT NOT NULL, -- 模型名称
	model_path TEXT NOT NULL, -- 模型路径，可能用户自定义路径
    api_domain TEXT NOT NULL, -- API访问域名
    api_service_id TEXT NOT NULL, -- API网关服务
    user_id TEXT DEFAULT NULL, -- 用户ID，平台服务为空
    power BIGINT NOT NULL DEFAULT 0, -- 部署的算力
    status TEXT DEFAULT 'none', -- 状态: none, downloading, enabled, disabled
    heartbeat_at TIMESTAMPTZ, -- 上次心跳时间
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* API网关服务表 */
DROP TABLE IF EXISTS api_services;
CREATE TABLE api_services (
    id TEXT PRIMARY KEY, -- 网关服务ID，资源调度返回的服务ID
    topo_id BIGINT NOT NULL, -- 所属拓扑域
    public_ip INET NOT NULL, -- 公网IP
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
    service_id TEXT NOT NULL, -- 模型服务ID
    request_limit BIGINT NOT NULL, -- 请求数限流（次/分钟）
    token_limit BIGINT NOT NULL,  -- Token限流（Tokens/分钟）
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (workspace_id, service_id)
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
    key_hash TEXT NOT NULL, -- 调用密钥
    user_id TEXT NOT NULL, -- 用户ID
    service_id BIGINT NOT NULL, -- 服务ID
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- 精确到毫秒
    status SMALLINT DEFAULT 0, -- 调用状态
    input_tokens BIGINT DEFAULT 0, -- 输入token数量
    output_tokens BIGINT DEFAULT 0, -- 输出token数量
	response_time_ms INT NOT NULL,  -- 响应耗时（毫秒）
    PRIMARY KEY (id, occurred_at)
);

SELECT create_hypertable('usage_logs', 'occurred_at');


