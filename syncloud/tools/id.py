import footprint
import uuid
import config


def getname(footprint):
    for name, f in config.footprints:
        if footprint.match(f):
            return name
    return None


def getmac():
    mac = uuid.getnode()
    mac_formated = ':'.join(("%012x" % mac)[i:i+2] for i in range(0, 12, 2))
    return mac_formated


def name():
    f = footprint.footprint()
    return getname(f)


class Id:
    def __init__(self, name, title, mac_address):
        self.name = name
        self.title = title
        self.mac_address = mac_address

def id():
    f = footprint.footprint()
    name = getname(f)
    if not name:
        raise Exception('Unknown footprint: {}'.format(f))
    mac_address = getmac()

    title = config.titles[name]
    id = Id(name, title, mac_address)
    return id