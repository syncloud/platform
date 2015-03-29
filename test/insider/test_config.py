import unittest

from convertible import to_json

from syncloud.insider.config import Port

class TestPortMapping(unittest.TestCase):

    def test_port_mapping(self):
        expected = '{"external_port": "8080", "local_port": "80"}'
        actual = to_json(Port("80", "8080"))
        self.assertEquals(expected, actual)