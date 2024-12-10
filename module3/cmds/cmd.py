class CMD:
    def __init__(self):        
        self.func_key = None        
        self.func = None
        self.tool = None
    
    def __call__(self):
        return self.func_key, self.func, self.tool