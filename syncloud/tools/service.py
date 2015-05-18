from syncloud.app import runner
from syncloud.app.logger import get_logger


class Service():
    def __init__(self):
        self.logger = get_logger('service')

    def stop(self, service):
        runner.call('service stop {}'.format(service), self.logger, shell=True)

    def start(self, service):
        runner.call('service start {}'.format(service), self.logger, shell=True)