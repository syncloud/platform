from sshcontrol import SshServer

SERVICE_NAME = 'ssh'
SERVICE_PROTOCOL = 'ssh'
SERVICE_TYPE = '_ssh._tcp'
SERVICE_PORT = 1022
SERVICE_URL = 'ssh'

class RemoteAccess:

    def __init__(self, insider):
        self.insider = insider
        self.ssh_server = SshServer('syncloud')

    def enable(self):
        private = self.ssh_server.setup(SERVICE_PORT, password_authentication=False)
        self.insider.add_service(SERVICE_NAME, SERVICE_PROTOCOL, SERVICE_TYPE, SERVICE_PORT, SERVICE_URL)
        return private

    def add_certificate(self):
        return self.ssh_server.add_certificate()

    def disable(self):
        self.insider.remove_service(SERVICE_NAME)
        self.ssh_server.remove()