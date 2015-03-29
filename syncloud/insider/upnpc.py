from syncloud.app import logger
from subprocess import check_output, CalledProcessError


class Upnpc:

    def __init__(self, local_ip):
        self.local_ip = local_ip
        self.logger = logger.get_logger('Upnpc')

    def external_ip(self):
        cmd = "upnpc -s | grep ExternalIPAddress | cut -d' ' -f3"
        self.logger.debug(cmd)
        try:
            return check_output(cmd, shell=True).strip()
        except CalledProcessError, e:
            self.logger.debug(e.output)
            raise e

    def mapped_external_ports(self, protocol):
        cmd = "upnpc -l | grep %s | awk '{ print $3 }' | cut -d'-' -f1" % protocol
        self.logger.debug(cmd)
        try:
            return map(int, check_output(cmd, shell=True).splitlines())
        except CalledProcessError, e:
            self.logger.debug(e.output)
            raise e

    def get_external_ports(self, protocol, local_port):
        cmd = "upnpc -l | grep %s | grep '%s:%s' | awk '{ print $3 }' | cut -d'-' -f1" % \
              (protocol, self.local_ip, local_port)
        self.logger.debug(cmd)
        try:
            return map(int, check_output(cmd, shell=True).splitlines())
        except CalledProcessError, e:
            self.logger.debug(e.output)
            raise e

    def remove(self, external_port):
        cmd = "upnpc -d {} TCP".format(external_port)
        self.logger.debug(cmd)
        try:
            check_output(cmd, shell=True)
        except CalledProcessError, e:
            self.logger.debug(e.output)
            raise e

    def add(self, local_port, external_port):
        cmd = "upnpc -a {} {} {} TCP".format(self.local_ip, local_port, external_port)
        self.logger.debug(cmd)
        try:
            return check_error(check_output(cmd, shell=True).strip())
        except CalledProcessError, e:
            self.logger.debug(e.output)
            raise e

    def port_open_on_router(self, ip, port):
        try:
            self.logger.debug(check_output('nc -w 1 {0} {1}'.format(ip, port), shell=True))
            self.logger.debug("{0}: port is taken".format(port))
            return True
        except CalledProcessError, e:
            self.logger.debug("{0}: port is available".format(port))
            self.logger.debug(e.output)
            return False


def check_error(output):
    if "failed" in output:
        error = next(iter([line for line in output.split('\n') if "failed" in line]))
        raise Exception('Unable to add mapping: ' + error)
    return output