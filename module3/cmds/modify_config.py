import json
from cmds.cmd import CMD

class CMDModifyConfig(CMD):
    def __init__(self):
        self.tool = {
            "type": "function",
            "function": {
                "name": "modify_config",
                "description": "修改服务配置",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "service_name": {
                            "type": "string",
                            "description": '服务名称，例如："nginx"',
                        },
                        "key": {
                            "type": "string",
                            "description": '配置的键值，例如："location"',
                        },
                        "value": {
                            "type": "string",
                            "description": '值 ，例如："/"',
                        },
                    },
                },
                "required": ["query_str", "key", "value"],
            },
        }
        self.func_key = 'modify_config'
        self.func = __class__.modify_config 
    
    @staticmethod
    def modify_config(service_name, key, value):
        print(f"modify config {service_name} {key} {value}")
        return json.dumps({"service": service_name, "key": key, "value": value})