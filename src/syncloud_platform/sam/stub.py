from subprocess import check_output
from syncloud_app import logger
from models import AppVersions
from os.path import join, isdir
import jsonpickle
import psutil
from shutil import rmtree, copytree

import convertible

SAM_BIN = '/opt/app/sam/bin/sam'

TEMP_SAM_PATH = '/tmp/sam-copy'
TEMP_SAM_BIN = join(TEMP_SAM_PATH, 'bin/sam')

class SamStub:

    def __init__(self, platform_config, info):
        self.info = info
        self.platform_config = platform_config
        self.logger = logger.get_logger('SamStub')

    def __get_sam_bin(self, app_id):
        sam_bin = SAM_BIN
        if app_id == 'sam':
            if isdir(TEMP_SAM_PATH):
                rmtree(TEMP_SAM_PATH, ignore_errors=True)
            copytree('/opt/app/sam', TEMP_SAM_PATH)
            sam_bin = TEMP_SAM_BIN
        return sam_bin

    def update(self, release=None):
        args = [SAM_BIN, 'update']
        if release:
            args += ['--release', release]
        return self.__run(args)

    def install(self, app_id):
        self.__run_detached('{0} install {1}'.format(SAM_BIN, app_id))

    def upgrade(self, app_id):
        sam_bin = self.__get_sam_bin(app_id)
        self.__run_detached('{0} upgrade {1}'.format(sam_bin, app_id))

    def remove(self, app_id):
        return self.__run([SAM_BIN, 'remove', app_id])

    def list(self):
        result = self.__run([SAM_BIN, 'list'])
        return [self._add_url(app_versions) for app_versions in convertible.to_object(result, convertible.List(item_type=AppVersions))]

    def _add_url(self, app_versions):
        app_versions.app.url = self.info.url(app_versions.app.id)
        return app_versions

    def user_apps(self):
        return [a for a in self.list() if not a.app.required]

    def installed_user_apps(self):
        return [a for a in self.user_apps() if a.installed_version]

    def installed_all_apps(self):
        return [a for a in self.list() if a.installed_version]

    def get_app(self, app_id):
        return next(a for a in self.list() if a.app.id == app_id)

    def __run_detached(self, command):
        # Think about adding twisted
        ssh_command = "ssh localhost -p {0} -o StrictHostKeyChecking=no 'nohup {1} </dev/null >/dev/null 2>&1 &'".format(
            self.platform_config.get_ssh_port(), command)
        self.logger.info('ssh command: {0}'.format(ssh_command))
        output = check_output(ssh_command, shell=True)
        self.logger.info(output)

    def __run(self, cmd_args):
        cmd_line = ' '.join(cmd_args)
        self.logger.info('cmd: {0}'.format(cmd_line))
        output = check_output(cmd_line, shell=True)
        result = jsonpickle.decode(output)
        return result['data']

    def is_running(self):
        for p in psutil.process_iter():
            for arg in p.cmdline():
                if SAM_BIN in arg or TEMP_SAM_PATH in arg:
                    return True
        return False
