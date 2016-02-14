from syncloud_app import logger


class Reader:

    def __init__(self, proc_cpuinfo='/proc/cpuinfo'):
        self.proc_cpuinfo = proc_cpuinfo
        self.log = logger.get_logger('Reader')
        
    def read(self):
        self.log.info('reading {}'.format(self.proc_cpuinfo))
        with open(self.proc_cpuinfo, 'r') as f:
            contents = f.read()
            # self.logger.debug('contents: {}'.format(contents))
            return contents