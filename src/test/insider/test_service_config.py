import unittest

from syncloud_platform.insider.config import Service
from test.insider.helpers import get_service_config


class TestServiceConfig(unittest.TestCase):

    def test_add_or_update(self):

        service_config = get_service_config([])

        service_config.add_or_update(Service("name", "proto", "type", 80, "url1"))
        service_config.add_or_update(Service("name", "proto", "type", 80, "url2"))
        service_config.add_or_update(Service("name1", "proto", "type", 81, "url3"))

        self.assertEquals(len(service_config.load()), 2)