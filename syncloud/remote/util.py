import os
from syncloud.app import runner


def copyfile(src, dst, log):
    runner.call('cp --preserve=mode,ownership,timestamps --force "{0}" "{1}"'.format(src, dst), log, shell=True)

def create_link(src, dst, log):
    runner.call('ln --symbolic --force "{0}" "{1}"'.format(src, dst), log, shell=True)

def rmfile(filename, ignore_errors=False):
    try:
        os.remove(filename)
    except Exception, ex:
        if not ignore_errors:
            raise ex

def parse_settings_line(line, delimeter):
    stripped = line.strip()
    if not stripped:
        return None
    if stripped.startswith('#'):
        return None
    index = stripped.find(delimeter)
    if index == -1:
        return None
    key = stripped[0:index]
    value = stripped[index+1: len(stripped)]
    return key.strip(), value.strip()


class Settings:
    def __init__(self, filename, delimeter):
        self.parameters = []
        self.delimeter = delimeter
        self.filename = filename
        self.__load()

    def __load(self):
        self.parameters = []
        if self.filename:
            f = open(self.filename, 'r')
            for line in f.readlines():
                setting = parse_settings_line(line, delimeter=self.delimeter)
                if setting:
                    self.parameters.append(setting)
            f.close()

    def save(self):
        f = open(self.filename, 'w')
        text = '\n'.join(['{0}{1}{2}'.format(key, self.delimeter, value) for (key, value) in self.parameters])
        f.write(text)
        f.close()

    def set(self, key, value, add_to_existing=False):
        if not add_to_existing:
            self.parameters = [(k, v) for (k, v) in self.parameters if k != key]
        self.parameters.append((key, value))
