import glob
from os.path import isfile
from subprocess import check_output


class Aggregator:
    def __init__(self, platform_config):
        self.platform_config = platform_config

    def get_logs(self):

        log_files = glob.glob(self.platform_config.get_log_sender_pattern())
        results = list(map(read_log, log_files))
        results.append(run('date'))
        results.append(run('dmesg'))
        results.append(run('mount'))
        results.append(run('journalctl -n 1000 --no-pager'))
        results.append(run('systemctl status --state=inactive snap.*'))
        results.append(run('COLUMNS=1000 top -n 1 -bc'))
        results.append(run('ping google.com -c 5'))
        results.append(run('df'))
        results.append(run('lsblk -o +UUID'))
        results.append(run('lsblk -Pp -o NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL'))
        results.append(run('ls -la /data'))
        results.append(run('uptime'))
        results.append(run('snap run platform.cli ipv4 public'))
        logs = '\n----------------------\n'.join(results)
        return logs


def read_log(filename):
    log = 'file: {0}\n\n'.format(filename)
    if isfile(filename):
        log += check_output('tail -100 {0}'.format(filename), shell=True).decode()
    else:
        log += '-- not found --'
    return log


def run(cmd):
    log = 'file: {0}\n\n'.format(cmd)
    log += check_output('{0} | tail -100 || true'.format(cmd), shell=True).decode()
    return log
