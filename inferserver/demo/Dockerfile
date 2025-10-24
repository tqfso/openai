FROM pytorch/pytorch:2.3.1-cuda12.1-cudnn8-devel

WORKDIR /app
COPY . /app
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install --upgrade triton

EXPOSE 8000
ENTRYPOINT ["uvicorn","main:app"]
CMD ["--host", "0.0.0.0", "--port", "8000"]