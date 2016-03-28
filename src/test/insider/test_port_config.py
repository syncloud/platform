import convertible

from syncloud_platform.insider.config import Port
from test.insider.helpers import get_port_config


def test_add_or_update():

    port_config = get_port_config([])

    port_config.add_or_update(Port(80, 10000, 'TCP'))
    port_config.add_or_update(Port(80, 10000, 'TCP'))
    port_config.add_or_update(Port(81, 10000, 'UCP'))
    port_config.add_or_update(Port(81, 10000, 'UDP'))

    assert len(port_config.load()) == 2


def test_pre_protocol_support():

    old_json = '[{"external_port": 80, "local_port": 80}]'

    port_config = get_port_config(convertible.from_json(old_json))

    assert port_config.get(80).protocol == 'TCP'
    port_config.add_or_update(Port(80, 10000, 'UDP'))
    assert port_config.get(80).protocol == 'UDP'

    print(open(port_config.filename, 'r').read())

    assert len(port_config.load()) == 1
