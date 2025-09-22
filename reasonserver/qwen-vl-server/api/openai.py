import time
import uuid
import requests
from io import BytesIO
from PIL import Image
from fastapi import APIRouter, Request
from fastapi import Depends
from pydantic import BaseModel
from typing import List, Union, Literal, Optional, Dict

# ---- FastAPI 路由和处理 ----

router = APIRouter()

# ---- 定义 OpenAI 风格请求结构 ----

class OpenAIImageURL(BaseModel):
    url: str

class OpenAITextContent(BaseModel):
    type: Literal["text"]
    text: str

class OpenAIImageContent(BaseModel):
    type: Literal["image_url"]
    image_url: OpenAIImageURL

OpenAIContentItem = Union[OpenAITextContent, OpenAIImageContent]

class OpenAIMessage(BaseModel):
    role: Literal["user", "assistant", "system"]
    content: Union[str, List[OpenAIContentItem]]

class OpenAIChatRequest(BaseModel):
    model: str
    messages: List[OpenAIMessage]
    max_tokens: Optional[int] = 128
    temperature: Optional[float] = 0.8

class OpenAIChatResponse(BaseModel):
    id: str
    object: str = "chat.completion"
    created: int
    model: str
    choices: List[Dict]
    usage: Optional[Dict] = None


def get_model(request: Request):
    return request.app.state.model

# ---- 转换 OpenAI -> 本地格式 ----
def convert_openai_to_internal(messages: List[OpenAIMessage]) -> List[Dict]:
    new_messages = []
    for msg in messages:
        contents = []
        if isinstance(msg.content, str):
            contents.append({"type": "text", "text": msg.content})
        else:
            for item in msg.content:
                if item.type == "text":
                    contents.append({"type": "text", "text": item.text})
                elif item.type == "image_url":
                    img_url = item.image_url.url
                    image = Image.open(BytesIO(requests.get(img_url).content)).convert("RGB")
                    image = image.resize((448, 448))  # 推荐尺寸，防OOM

                    contents.append({"type": "image", "image": image})
        new_messages.append({
            "role": msg.role,
            "content": contents
        })
    return new_messages

@router.post("/chat/completions", response_model=OpenAIChatResponse)
def chat_completions(req: OpenAIChatRequest, request: Request):

    model = get_model(request)
    model_inputs = {
        "messages": convert_openai_to_internal(req.messages),
        "temperature": req.temperature,
        "max_tokens": req.max_tokens,
    }

    result = model.generate(model_inputs)
    output_text = result["output"]
    usage = result["usage"]

    return OpenAIChatResponse(
        id="chatcmpl-" + str(uuid.uuid4()),
        created=int(time.time()),
        model=req.model,
        choices=[
            {
                "index": 0,
                "message": {
                    "role": "assistant",
                    "content": output_text
                },
                "finish_reason": "stop"
            }
        ],
        usage=usage
    )