import syncloud_platform.importlib

import unittest

from syncloud_platform.insider.config import Port
from test.insider.helpers import get_port_config


class TestPortConfig(unittest.TestCase):

    def test_add_or_update(self):

        port_config = get_port_config([])

        port_config.add_or_update(Port(80, 10000))
        port_config.add_or_update(Port(80, 10000))
        port_config.add_or_update(Port(81, 10000))
        port_config.add_or_update(Port(81, 10000))

        self.assertEquals(len(port_config.load()), 2)