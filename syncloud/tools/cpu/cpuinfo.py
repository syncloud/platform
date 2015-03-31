import re
from syncloud.app import logger


class CpuInfo:

    def __init__(self, reader):
        self.reader = reader
        self.logger = logger.get_logger('CpuInfo')

    def vendor_id(self):
        return self.value('vendor_id')

    def hardware(self):
        return self.value('hardware')

    def value(self, key):
        lines = self.reader.read().splitlines()
        for line in lines:
            result = self.parse(line)
            if result:
                k, v = result
                # self.logger.debug('{} = {}'.format(k, v))
                if k.lower() == key.lower():
                    return v

    def parse(self, line):
        m = re.match('(.*):(.*)', line)
        if m:
            key = m.group(1).strip()
            value = m.group(2).strip()
            return key, value

