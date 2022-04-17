from syncloud_platform.insider.config import Port
from syncloud_platform.insider.device_info import DeviceInfo
from test.insider.helpers import get_user_platform_config, get_port_config


def test_url_activated_free():

    user_platform_config = get_user_platform_config()
    user_platform_config.update_domain('device', 'token')
    user_platform_config.update_redirect('syncloud.it')
    user_platform_config.set_manual_access_port(10000)
    user_platform_config.set_redirect_enabled(True)

    device_info = DeviceInfo(user_platform_config, port_config)

    assert device_info.url('app') == 'https://app.device.syncloud.it:10000'


def test_url_activated_custom():

    user_platform_config = get_user_platform_config()
    user_platform_config.set_custom_domain('example.com')
    user_platform_config.update_redirect('syncloud.it')
    user_platform_config.set_manual_access_port(10000)
    user_platform_config.set_redirect_enabled(False)

    device_info = DeviceInfo(user_platform_config, port_config)

    assert device_info.url('app') == 'https://app.example.com:10000'


def test_url_non_activated():

    user_platform_config = get_user_platform_config()
    user_platform_config.update_domain('device', 'token')
    user_platform_config.update_redirect('syncloud.it')
    user_platform_config.set_manual_access_port(10000)
    user_platform_config.set_redirect_enabled(False)

    device_info = DeviceInfo(user_platform_config, port_config)

    assert device_info.url('app') == 'https://app.localhost:10000'
