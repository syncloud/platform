import id
import footprint
import network
from syncloud.tools import env


class Facade:

    def id(self):
        return id.id()

    def name(self):
        return id.name()

    def footprint(self):
        return footprint.footprint()

    def local_ip(self):
        return network.local_ip()

    def usr_local_dir(self):
        return env.usr_local_dir()

    def root_dir_prefix(self):
        return env.root_dir_prefix()