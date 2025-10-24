FROM vllm/vllm-openai:latest

# 解决提示词为音频问题
RUN pip install librosa soundfile

# 设置环境变量，解决FlashInfer异常问题
ENV VLLM_USE_FLASHINFER_SAMPLER=0