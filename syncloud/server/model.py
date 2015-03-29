class Credentials:
    def __init__(self, login, password, key):
        self.login = login
        self.password = password
        self.key = key


class Site():
    def __init__(self, endpoint):
        service = endpoint.service
        self.name = service.name.lower()
        self.url = "{0}://{1}:{2}/{3}".format(service.protocol, endpoint.external_host, endpoint.external_port, service.url)