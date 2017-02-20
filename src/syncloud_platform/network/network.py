import netifaces
from syncloud_app import logger

class Network:
    def __init__(self):
        self.log = logger.get_logger('network')

    def interfaces(self):
        ifaces=netifaces.interfaces()
        return [self.__convert(iface) for iface in ifaces]

    def __convert(self, iface):
        addrs=netifaces.ifaddresses(iface)
        ipv4=addrs[netifaces.AF_INET]
        ipv6=addrs[netifaces.AF_INET6]
        return dict(iface=iface, ipv4=ipv4, ipv6=ipv6)

