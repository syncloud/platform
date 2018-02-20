import glob
from os.path import isfile
from subprocess import check_output


class Aggregator:
    def __init__(self, platform_config):
        self.platform_config = platform_config

    def get_logs(self):

        log_files = glob.glob(self.platform_config.get_log_sender_pattern())
        log_files.append('/var/log/sam.log')
        log_files.append('/var/log/syslog')
        results  = map(read_log, log_files)
        results.append(run('dmesg'))
        results.append(run('mount'))
        results.append(run('journalctl'))
        results.append(run('df'))
        results.append(run('lsblk'))
        logs = '\n----------------------\n'.join(results)
        return logs


def read_log(filename):
    log = 'file: {0}\n\n'.format(filename)
    if isfile(filename):
        log += check_output('tail -100 {0}'.format(filename), shell=True)
    else:
        log += '-- not found --'
    return log


def run(cmd):
    log = 'file: {0}\n\n'.format(cmd)
    log += check_output('{0} | tail -100 || true'.format(cmd), shell=True)
    return log

