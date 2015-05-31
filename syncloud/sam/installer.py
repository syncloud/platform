import os
from os.path import isfile
from subprocess import check_output
from syncloud.app import logger
from syncloud.sam.manager import get_sam
import wget
import tarfile


ROOT_PATH = '/opt/app'

class Installer:
    def __init__(self):
        self.log = logger.get_logger('sam.installer')
        self.sam = get_sam()

    def install(self, app, from_file=None):
        if not from_file:
            archive = '{0}.tar.gz'.format(app)
            arch = check_output('uname -m', shell=True).strip()
            url = 'http://apps.syncloud.org/{0}/{1}/{2}'.format(self.sam.get_release(), arch, archive)

            self.log.info('installing from: {0}'.format(url))

            temp = '/tmp/{0}'.format(archive)
            if isfile(temp):
                os.remove(temp)

            self.log.info("saving {0} to {1}".format(url, temp))
            filename = wget.download(url, temp)
        else:
            filename = from_file
            self.log.info('installing from: {0}'.format(from_file))

        self.log.info("extracting {0}".format(filename))
        tarfile.open(filename).extractall(ROOT_PATH)
