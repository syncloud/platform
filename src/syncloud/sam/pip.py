from syncloud.app import runner
from syncloud.app import logger
import logging

import requests
import json

class Pip:

    def __init__(self, pypi_index, raise_on_error=True):
        self.logger = logger.get_logger('Pip')
        self.pypi_index = pypi_index
        self.raise_on_error = raise_on_error

    def install(self, name, version):
        if self.pypi_index:
            cmd_args = ['pip2', 'install', '--no-binary', ':all:', '--index-url', self.pypi_index, '{}=={}'.format(name, version)]
        else:
            cmd_args = ['pip2', 'install', '--no-binary', ':all:', '{}=={}'.format(name, version)]
        exit_code = runner.call(' '.join(cmd_args), self.logger, stdout_log_level=logging.INFO, shell=True)
        if self.raise_on_error and not exit_code == 0:
            raise Exception('failed command: '+' '.join(cmd_args))

    def uninstall(self, name):
        cmd_args = ['yes', '|', 'pip2', 'uninstall', name]
        exit_code = runner.call(' '.join(cmd_args), self.logger, stdout_log_level=logging.INFO, shell=True)
        if self.raise_on_error and not exit_code == 0:
            raise Exception('failed command: '+' '.join(cmd_args))

    def last_version(self, name):
        url = "https://pypi.python.org/pypi/%s/json" % (name,)
        response = requests.get(url)
        data = json.loads(response.content)
        version = data['info']['version']
        return str(version)

    def log_version(self, name):
        cmd_args = ['pip2', 'freeze', '|', 'grep', name]
        exit_code = runner.call(' '.join(cmd_args), self.logger, stdout_log_level=logging.INFO, shell=True)
        if self.raise_on_error and not exit_code == 0:
            raise Exception('failed command: '+' '.join(cmd_args))
