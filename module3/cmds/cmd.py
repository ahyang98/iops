class CMD:
    def __init__(self):
        self.tool = None
        self.func_key = None
        self.func = None
    
    def __call__(self):
        return self.func_key, self.tool, self.func