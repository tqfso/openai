FROM reasoning/cosyvoice:0.1.1

# 设置默认工作目录
WORKDIR /workspace

# 拷贝可执行文件
COPY server.py /workspace/CosyVoice/server.py

# 暴露端口
EXPOSE 8000

# 启动入口
ENTRYPOINT [ "python", "/workspace/CosyVoice/server.py" ]

