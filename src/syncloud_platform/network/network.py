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
        ipv4=self.__addr(addrs[netifaces.AF_INET])
        self.log.info(addrs[netifaces.AF_INET])
        ipv6=self.__addr(addrs[netifaces.AF_INET6])
        self.log.info(addrs[netifaces.AF_INET6])
        return dict(iface=iface, ipv4=ipv4, ipv6=ipv6)

    def __addr(self, info):
        if 'addr' in info:
            return info['addr']
        return None

