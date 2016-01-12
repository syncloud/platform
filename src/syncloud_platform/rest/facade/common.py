from os import listdir
from os.path import isfile, join
from subprocess import check_output

html_prefix = '/server/html'
rest_prefix = '/server/rest'


class Common:
    def __init__(self, platform_config):
        self.platform_config = platform_config
        self.log_root = self.platform_config.get_log_root()

    def get_logs(self):

        log_files = [join(self.log_root, f) for f in listdir(self.log_root) if isfile(join(self.log_root, f))]
        log_files.append('/var/log/sam.log')

        logs = '\n----------------------\n'.join(map(self.read_log, log_files))

        return logs

    def read_log(self, filename):
        log = 'file: {0}\n\n'.format(filename)
        if isfile(filename):
            log += check_output('tail -100 {0}'.format(filename), shell=True)
        else:
            log += '-- not found --'
        return log

