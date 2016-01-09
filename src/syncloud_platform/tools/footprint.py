import psutil
from syncloud_platform.tools.cpu.cpuinfo import CpuInfo
from syncloud_platform.tools.cpu.reader import Reader

from subprocess import check_output


def match_contains(pattern, value):
    if pattern is None:
        return True
    if value is None:
        return False
    return pattern in value


class Footprint:
    def __init__(self, cpu_hardware=None, cpu_count=None, mem_size=None, vendor_id=None, lsusb=None):
        self.cpu_hardware = cpu_hardware
        self.cpu_count = cpu_count
        self.mem_size = mem_size
        self.vendor_id = vendor_id
        self.lsusb = lsusb

    def match(self, pattern):
        this = self.__dict__
        for (k, v) in pattern.__dict__.iteritems():
            if not self.match_member(k, v, this[k]):
                return False
        return True

    def match_member(self, name, pattern, value):
        if name == 'lsusb':
            return match_contains(pattern, value)
        else:
            return pattern is None or value == pattern

    def __str__(self):
        return str(self.__dict__)


def lsusb():
    return check_output(['lsusb'])


def footprint():
    cpu_count = psutil.cpu_count()
    mem_size = psutil.virtual_memory().total
    info = CpuInfo(Reader())
    cpu_hardware = info.hardware()
    vendor_id = info.vendor_id()
    return Footprint(cpu_hardware, cpu_count, mem_size, vendor_id, None)