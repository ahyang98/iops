import json
from cmds.cmd import CMD

class CMDRestartService(CMD):
    def __init__(self):
        self.tool = {
            "type": "function",
            "function": {
                "name": "restart",
                "description": "重启服务",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "service_name": {
                            "type": "string",
                            "description": '服务名称，例如："nginx"',
                        },
                    },
                },
                "required": ["service_name"],
            },
        }
        self.func_key = 'restart'
        self.func = self.__class__.restart 
    
    @staticmethod
    def restart(service_name):
        print(f"restart service {service_name}")
        return json.dumps({"service": service_name})