from subprocess import check_output
from models import AppVersions
import jsonpickle

import convertible

class SamStub:

    def __init__(self):
        self.sam = '/opt/app/sam/bin/sam'

    def update(self, release=None):
        args = [self.sam, 'update']
        if release:
            args += ['--release', release]
        return self.__run(args)

    def install(self, app_id):
        return self.__run([self.sam, 'install', app_id])

    def remove(self, app_id):
        return self.__run([self.sam, 'remove', app_id])

    def list(self):
        result = self.__run([self.sam, 'list'])
        return convertible.to_object(result, convertible.List(item_type=AppVersions))

    def __run(self, cmd_args):
        cmd_line = ' '.join(cmd_args)
        output = check_output(cmd_line, shell=True)
        result = jsonpickle.decode(output)
        return result['data']