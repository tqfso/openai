import json
from openai import OpenAI
from config import BASE_URL

# 测试获取天气信息

model_name = "Qwen/Qwen3-1.7B"

client = OpenAI(
    base_url=BASE_URL,
    api_key="EMPTY"
)

tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "获取指定城市的当前天气信息",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {
                        "type": "string",
                        "description": "城市名称，例如：北京、上海、深圳"
                    }
                },
                "required": ["city"]
            }
        }
    }
]

def get_weather(city: str) -> str:
    """
    模拟获取天气信息的函数。
    在实际应用中，这里应该调用真实的天气 API（如和风天气、OpenWeatherMap 等）。
    """
    # 模拟数据 - 仅用于演示
    weather_data = {
        "深圳": {"temperature": "28°C", "condition": "晴", "humidity": "65%"},
        "北京": {"temperature": "22°C", "condition": "多云", "humidity": "45%"},
        "上海": {"temperature": "25°C", "condition": "小雨", "humidity": "80%"}
    }
    
    if city in weather_data:
        return json.dumps(weather_data[city], ensure_ascii=False)
    else:
        return json.dumps({"error": f"未找到城市 {city} 的天气数据"}, ensure_ascii=False)

messages = [
    {"role": "system", "content": "你是一名气象助手"},
    {"role": "user", "content": "今天深圳的天气怎么样？"}
]

print("用户提问:", messages[1]["content"])

response = client.chat.completions.create(
    model=model_name,
    messages=messages,
    tools=tools,
    tool_choice="auto",
    max_tokens=256,
    temperature=0.2,
    top_p=0.95,
)

response_message = response.choices[0].message

if hasattr(response_message, 'tool_calls') and response_message.tool_calls:
    # 模型希望调用工具
    tool_call = response_message.tool_calls[0]
    
    if tool_call.function.name == "get_weather":
        # 解析参数
        function_args = json.loads(tool_call.function.arguments)
        city = function_args.get("city")
        
        print(f"模型决定调用工具: get_weather(city='{city}')")
        
        # 执行工具函数
        try:
            tool_response = get_weather(city)
            print(f"工具返回结果: {tool_response}")
        except Exception as e:
            tool_response = json.dumps({"error": f"调用工具时出错: {str(e)}"})
            print(f"工具调用出错: {tool_response}")
        
        # 将工具执行结果添加到消息历史，并再次调用模型生成最终回复
        messages.append(response_message)
        messages.append({
            "tool_call_id": tool_call.id,
            "role": "tool",
            "name": "get_weather",
            "content": tool_response,
        })
        
        # 再次调用模型，让其基于工具结果生成最终回答
        final_response = client.chat.completions.create(
            model=model_name,
            messages=messages,
            max_tokens=256,
            temperature=0.6,
            top_p=0.95,
            stream=True,
        )

        # 实时打印每一段内容
        final_answer = ""
        for chunk in final_response:
            delta = chunk.choices[0].delta
            if delta.content:
                print(delta.content, end='', flush=True)
                final_answer += delta.content
        
        # 非流式打印
        # final_answer = final_response.choices[0].message.content
        # print("模型最终回复:", final_answer)
        
    else:
        print("未知的工具调用:", tool_call.function.name)
else:
    # 模型没有调用工具，直接生成了回复
    direct_answer = response_message.content
    print("模型直接回复:", direct_answer)