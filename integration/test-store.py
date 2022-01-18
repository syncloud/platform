import os

import pytest
from subprocess import check_output
from syncloudlib.integration.hosts import add_host_alias

TMP_DIR = '/tmp/syncloud'


@pytest.fixture(scope="session")
def module_setup(request, device, artifact_dir):
    def module_teardown():
        device.run_ssh('journalctl > {0}/store-test.journalctl.log'.format(TMP_DIR), throw=False)
        check_output('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def test_start(module_setup, app, device_host, domain):
    add_host_alias(app, device_host, domain)


def test_sore(device):
    channel = os.environ["DRONE_BRANCH"]
    if channel == 'stable':
        channel = 'rc'
    device.activated()
    device.run_ssh('snap refresh platform --channel={0} --amend'.format(channel))
