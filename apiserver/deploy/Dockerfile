FROM openai/ubuntu:22.04

# 设置默认工作目录
WORKDIR /app

# 拷贝配置文件
RUN mkdir -p /app/config
COPY config/config.yaml /app/config/config.yaml

# 拷贝可执行文件
COPY out/apiserver /app/
RUN chmod +x /app/apiserver

# 暴露端口
EXPOSE 8443

# 启动入口
ENTRYPOINT [ "/app/apiserver" ]

