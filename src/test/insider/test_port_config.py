import logging

from syncloud_app import logger

from syncloud_platform.insider.config import Port
from test.insider.helpers import get_port_config

logger.init(level=logging.DEBUG, console=True)


def test_add_or_update():

    port_config = get_port_config([])

    port_config.add_or_update(Port(80, 10000))
    port_config.add_or_update(Port(80, 10000))
    port_config.add_or_update(Port(81, 10000))
    port_config.add_or_update(Port(81, 10000))

    assert len(port_config.load()) == 2