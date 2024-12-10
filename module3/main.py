from function_call import FunctionCall
from cmds.apply_manifest import CMDApplyManifest
from cmds.modify_config import CMDModifyConfig
from cmds.restart import CMDRestartService


def main():
    function_call = FunctionCall()
    function_call.register(CMDApplyManifest()).\
    register(CMDModifyConfig()).\
    register(CMDRestartService())
    while True:
        resp = function_call.conversation()
        print("LLM Res: ", resp)
    

if __name__ == "__main__":    
    main()