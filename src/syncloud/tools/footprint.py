import psutil
from syncloud.tools.cpu.cpuinfo import CpuInfo
from syncloud.tools.cpu.reader import Reader


class Footprint:
    def __init__(self, cpu_hardware=None, cpu_count=None, mem_size=None, vendor_id=None):
        self.cpu_hardware = cpu_hardware
        self.cpu_count = cpu_count
        self.mem_size = mem_size
        self.vendor_id = vendor_id

    def match(self, pattern):
        this = self.__dict__
        for (k, v) in pattern.__dict__.iteritems():
            if v and not this[k] == v:
                return False
        return True

    def __str__(self):
        return str(self.__dict__)


def footprint():
    cpu_count = psutil.cpu_count()
    mem_size = psutil.virtual_memory().total
    info = CpuInfo(Reader())
    cpu_hardware = info.hardware()
    vendor_id = info.vendor_id()
    return Footprint(cpu_hardware, cpu_count, mem_size, vendor_id)