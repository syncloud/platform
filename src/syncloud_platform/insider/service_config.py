from os.path import join, isfile
from os import remove

import convertible

from syncloud_platform.config.config import PLATFORM_APP_NAME
from syncloud_platform.tools.app import get_app_data_root

SERVICE_CONFIG_NAME = 'services.json'


class ServiceConfig:

    def __init__(self, config_dir):
        self.filename = join(config_dir, SERVICE_CONFIG_NAME)

    def load(self):
        items = convertible.read_json(self.filename)
        if not items:
            return []
        return items

    def save(self, items):
        convertible.write_json(self.filename, items)

    def add_or_update(self, mapping):
        if self.get(mapping.name):
            self.__update(mapping)
        else:
            self.__add(mapping)

    def __add(self, item):
        items = self.load()
        items.append(item)
        self.save(items)

    def __update(self, new_service):
        service_list = self.load()
        mapping = next((m for m in service_list if m.name == new_service.name), None)
        loc = service_list.index(mapping)
        service_list[loc] = new_service
        self.save(service_list)

    def get(self, name):
        items = self.load()
        item = next((m for m in items if m.name == name), None)
        return item

    def get_by_port(self, local_port):
        items = self.load()
        item = next((m for m in items if m.port == local_port), None)
        return item

    def remove(self, name):
        items = self.load()
        new_items = [i for i in items if i.name != name]
        self.save(new_items)

    def remove_all(self):
        if isfile(self.filename):
            remove(self.filename)
