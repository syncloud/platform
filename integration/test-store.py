import os

from syncloudlib.integration.hosts import add_host_alias


def test_start(app, device_host, domain):
    add_host_alias(app, device_host, domain)


def test_sore(device):
    channel = os.environ["DRONE_BRANCH"]
    if channel == 'stable':
        channel = 'rc'
    device.login()
    device.run_ssh('snap refresh platform --channel={0}'.format(channel))
