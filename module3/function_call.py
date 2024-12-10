from openai import OpenAI
import json
from .cmds.cmd import CMD

class FunctionCall:
    def __init__(self):
        self.client = OpenAI()
        self.system = {
            "role":"system",
            "content": "你是一个IAC助手，可以修改配置，操作服务",
            }
        self.tools = []
        self.funcs = {}
    
    def register(self, cmd: CMD):
        if cmd.key and cmd.func:
            self.funcs[cmd.key]=cmd.func

        if cmd.tool:
            self.tools.append(cmd.tool)
    
    def conversation(self):
        query = input("请输入命令: ")
        messages = []
        messages.append(self.system)
        messages.append({
            "role": "user",
            "content": query
        })

        response = self.client.chat.completions.create(
            model="gpt-4o",
            messages=messages,
            tools=self.tools,
            tool_choice="auto"
        )        
        response_message = response.choices[0].message
        tool_calls = response_message.tool_calls
        
        print("\nChatGPT want to call function: ", tool_calls)
        if tool_calls is None:
            print("not tool_calls")
            return
        
        messages.append(response_message)

        for tool_call in tool_calls:
            func_name = tool_call.function.name
            func_to_call = self.funcs[func_name]
            args = json.loads(tool_call.function.arguments)
            func_resp = func_to_call(**args)
            messages.append(
                {
                    "tool_call_id": tool_call.id,
                    "role": "tool",
                    "name": func_name,
                    "content": func_resp,
                }
            )
        
        resp = self.client.completions.create(
            model="gpt-4o",
            messages=messages
        )
        return resp.choices[0].message.content

