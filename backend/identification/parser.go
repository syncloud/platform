package identification

import (
	"net"
)

func GetMac() (string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, ifa := range ifas {
		addr := ifa.HardwareAddr.String()
		if len(ifa.HardwareAddr) >= 6 && ifa.Name != "" {
			return addr, nil
		}
	}
	return "", nil

}

/*
class IdConfig:
def __init__(self, filename):
self.parser = ConfigParser()
if isfile(filename):
self.parser.read(filename)

def __get(self, key, default=None):
if self.parser.has_option('id', key):
return self.parser.get('id', key)
return default

def name(self):
return self.__get('name', default='unknown')

def title(self):
return self.__get('title', default='Unknown')


class Id:
def __init__(self, name, title, mac_address):
self.name = name
self.title = title
self.mac_address = mac_address


def id(id_config_filename='/etc/syncloud/id.cfg'):
id_config = IdConfig(id_config_filename)
mac_address = getmac()
name = id_config.name()
title = id_config.title()
id = Id(name, title, mac_address)
return id

*/
