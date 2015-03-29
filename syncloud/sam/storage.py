from os.path import isdir, dirname, isfile
from os import makedirs

import convertible
from syncloud.app import logger
import models


class Applications:

    def __init__(self, filename):
        self.filename = filename

    def list(self):
        apps_data = convertible.read_json(self.filename, models.Apps)
        if not apps_data:
            return []
        return apps_data.apps


class Versions:

    def __init__(self, filename, allow_latest=False):
        self.filename = filename
        self.__create_folder(filename)
        self.allow_latest = allow_latest
        self.logger = logger.get_logger('versions')


    def __create_folder(self, filename):
        path = dirname(filename)
        if not isdir(path):
            makedirs(path)

    def __read(self):
        versions = {}
        if isfile(self.filename):
            with open(self.filename, 'r') as f:
                lines = f.readlines()
                for line in lines:
                    if '=' in line:
                        id, version = line.strip().split('=', 1)
                        versions[id.strip()] = version.strip()
                    else:
                        if self.allow_latest:
                            raise Exception('')
                        else:
                            id = line
                            versions[id.strip()] = None
        return versions

    def __write(self, versions):
        with open(self.filename, 'w+') as f:
            for id, version in versions.iteritems():
                f.write('{}={}\n'.format(id, version))

    def version(self, name):
        # self.logger.debug('getting versions for {}'.format(name))
        versions = self.__read()
        # self.logger.debug('found {} versions'.format(len(versions)))
        return versions.get(name, None)

    def remove(self, name):
        versions = self.__read()
        del versions[name]
        self.__write(versions)

    def update(self, name, version):
        versions = self.__read()
        versions[name] = version
        self.__write(versions)
