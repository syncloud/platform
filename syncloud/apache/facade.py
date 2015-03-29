import site
from syncloud.apache.setup import Setup
from syncloud.apache.system import System


class ApacheFacade:

    def __init__(self):
        self.system = System()
        self.setup = Setup(self.system)

    def activate(self, hostname):
        ports = self.setup.activate(hostname)
        self.system.restart()
        return ports

    def add_http_site(self, name, config_file):
        return site.add(name, config_file, False)

    def add_https_site(self, name, config_file):
        return site.add(name, config_file, True)

    def add_site(self, protocol, name, config_file):
        return site.add(name, config_file, protocol == 'https')

    def remove_http_site(self, name):
        site.remove(name, False)

    def remove_https_site(self, name):
        site.remove(name, True)

    def restart(self):
        self.system.restart()
