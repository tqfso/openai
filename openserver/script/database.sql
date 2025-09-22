
/* 网关表 */
DROP TABLE IF EXISTS api_gateways;
CREATE TABLE api_gateways (
    id BIGSERIAL PRIMARY KEY,
    type SMALLINT NOT NULL DEFAULT 0, -- 网关类型
    service_id TEXT, -- 网关服务ID
    topo_id BIGINT NOT NULL, -- 所属拓扑域
    public_ip INET NOT NULL, -- 公网IP
    access_key TEXT NOT NULL, -- 访问密钥,存储HASH值
    status SMALLINT DEFAULT 0, -- 状态
    heartbeat_at TIMESTAMPTZ, -- 上次心跳时间
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()

);

/* 预置模型表*/
DROP TABLE IF EXISTS platform_models;
CREATE TABLE platform_models (
    id BIGSERIAL PRIMARY KEY,
    model_name TEXT UNIQUE NOT NULL, -- 模型名称如 Qwen/Qwen3-Reranker-8B
    provider TEXT, -- 深度求索、通义实验室等
    languages TEXT[], -- 支持语言 ['zh', 'en']
    classes TEXT[] NOT NULL, -- 文本生成/图片生成/语音识别等
    extended_ability TEXT[], -- 扩展能力如: ['function', 'mcp', 'reasoning', 'batch']
    context_length BIGINT NOT NULL, -- 最大上下文长度
	deploy_info JSONB, -- 部署信息包括：推理镜像列表 => 运行命令、运行参数、可用加速卡等
    finetune_info JSONB, -- 微调训练信息
    status SMALLINT DEFAULT 0, -- 状态
    description TEXT, -- 描述
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

/* 模型服务表 */
DROP TABLE IF EXISTS model_services;
CREATE TABLE model_services (
    id BIGSERIAL PRIMARY KEY,
	topo_id BIGINT NOT NULL, -- 所属拓扑域
    service_id TEXT, -- 模型服务ID
	service_name TEXT NOT NULL, -- 服务名称
    model_name TEXT NOT NULL, -- 模型名称
	model_path TEXT NOT NULL, -- 模型路径
    user_id TEXT DEFAULT NULL, -- 用户ID    
	is_platform BOOLEAN GENERATED ALWAYS AS ((user_id IS NULL)) STORED,
    power BIGINT NOT NULL DEFAULT 0, -- 部署的算力
    status SMALLINT DEFAULT 0, -- 状态
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()

);

/* API密钥表 */
DROP TABLE IF EXISTS api_keys;
CREATE TABLE api_keys (
    id BIGSERIAL PRIMARY KEY,
    key_hash TEXT NOT NULL UNIQUE, -- 密钥HASH
    user_id TEXT NOT NULL, -- 用户ID
    permissions JSONB, -- 权限控制：{ "open_models": ["qwen", "deepseek"], "private_models": ["res-123"] }
    status SMALLINT DEFAULT 0, -- 状态
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


