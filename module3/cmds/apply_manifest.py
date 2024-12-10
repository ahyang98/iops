import json
from cmd import CMD

class CMDApplyManifest(CMD):
    def __init__(self):
        self.tool = {
            "type": "function",
            "function": {
                "name": "apply_manifest",
                "description": "重启服务",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "resource_type": {
                            "type": "string",
                            "description": '资源对象，例如："deployment"',
                        },
                        "image": {
                            "type": "string",
                            "description": '资源对象，例如："nginx:1.0"',
                        },
                    },
                },
                "required": ["resource_type", "image"],
            },
        }
        self.func_key = 'apply_manifest'
        self.func = __class__.apply_manifest 
    
    @staticmethod
    def apply_manifest(resource_type，image):
        print(f"apply manifest resource {resource_type} image {image}")
        return json.dumps({"resource": resource_type，"image": image})