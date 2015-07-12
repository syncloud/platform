from string import join
from syncloud_app import runner
from syncloud_app.logger import get_logger


class Apt():

    def __init__(self):
        self.logger = get_logger('apache.apt')

    def install(self, apps):
        runner.call('apt-get -y install {}'.format(join(apps, " ")), self.logger, shell=True)

    def update(self):
        runner.call('apt-get update', self.logger, shell=True)

    def add_repo(self, repo):
        runner.call("add-apt-repository -y {0}".format(repo), self.logger, shell=True)
