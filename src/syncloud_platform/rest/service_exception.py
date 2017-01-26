class ServiceException(Exception):
    def __init__(self, message):
        Exception.__init__(self)
        self.message = message