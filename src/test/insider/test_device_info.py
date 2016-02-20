import logging

from syncloud_app import logger

from syncloud_platform.insider.config import Port
from syncloud_platform.insider.device_info import DeviceInfo
from test.insider.helpers import get_user_platform_config, get_port_config


logger.init(level=logging.DEBUG, console=True)


def test_url_with_external_access():

    user_platform_config = get_user_platform_config()
    user_platform_config.update_domain('device', 'token')
    user_platform_config.update_redirect('syncloud.it', 'api.url')
    user_platform_config.update_device_access(True, 'http')

    port_config = get_port_config([Port(80, 10000)])

    device_info = DeviceInfo(user_platform_config, port_config)

    assert device_info.url('app') == 'http://app.device.syncloud.it:10000'
