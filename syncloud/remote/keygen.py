import os
import tempfile
from os import path
from os.path import exists
from syncloud.app import logger

from syncloud.app import runner


class KeyGen:

    def __init__(self):
        self.logger = logger.get_logger('KeyGen')

    def generate(self, type, bits=None):
        key_file = path.join(tempfile.mkdtemp(), 'key')
        self.generate_into_file(type, key_file, bits)
        return self.read(key_file), self.read(key_file + '.pub')

    def generate_into_file(self, type, key_file, bits=None, overwrite=False):
        if overwrite:
            self.clean(key_file)
            self.clean(key_file + '.pub')
        command = "ssh-keygen -t {} -f {} -N ''".format(type, key_file)
        if bits:
            command + " -b {}".format(str(bits))
        runner.call(command, self.logger, shell=True)

    def read(self, key_file):
        with open(key_file, 'r') as f:
            return f.read()

    def clean(self, key_file):
        if exists(key_file):
            os.remove(key_file)