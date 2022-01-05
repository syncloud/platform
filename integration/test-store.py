from syncloudlib.integration.hosts import add_host_alias


def test_start(app, device_host, domain):
    add_host_alias(app, device_host, domain)


def test_sore(device):
    device.run_ssh('snap install platform --channel=rc')
    channel=rc')
