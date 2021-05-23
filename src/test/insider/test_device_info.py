from syncloud_platform.insider.config import Port
from syncloud_platform.insider.device_info import DeviceInfo
from test.insider.helpers import get_user_platform_config, get_port_config


def test_url_activated_free():

    user_platform_config = get_user_platform_config()
    user_platform_config.set_activated()
    user_platform_config.update_domain('device', 'token')
    user_platform_config.update_redirect('syncloud.it', 'api.url')
    user_platform_config.update_device_access(False, True, '1.1.1.1', 80, 443)
    user_platform_config.set_redirect_enabled(True)

    port_config = get_port_config([Port(443, 10000, 'TCP')])

    device_info = DeviceInfo(user_platform_config, port_config)

    assert device_info.url('app') == 'https://app.device.syncloud.it:10000'


def test_url_non_activated():

    user_platform_config = get_user_platform_config()
    user_platform_config.update_domain('device', 'token')
    user_platform_config.update_redirect('syncloud.it', 'api.url')
    user_platform_config.update_device_access(False, True, '1.1.1.1', 80, 443)
    user_platform_config.set_redirect_enabled(True)

    port_config = get_port_config([Port(443, 10000, 'TCP')])

    device_info = DeviceInfo(user_platform_config, port_config)

    assert device_info.url('app') == 'https://app.localhost:10000'
