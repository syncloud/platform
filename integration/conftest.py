import pytest

SYNCLOUD_INFO = 'syncloud.info'


def pytest_addoption(parser):
    parser.addoption("--domain", action="store")
    parser.addoption("--device-host", action="store")
    parser.addoption("--app-archive-path", action="store")


@pytest.fixture(scope='session')
def app_archive_path(request):
    return request.config.getoption("--app-archive-path")


@pytest.fixture(scope='session')
def device_host(request):
    return request.config.getoption("--device-host")


@pytest.fixture(scope='session')
def domain(request):
    return request.config.getoption("--domain")


@pytest.fixture(scope='session')
def main_domain():
    return SYNCLOUD_INFO


@pytest.fixture(scope='session')
def device_domain(domain, main_domain):
    return '{0}.{1}'.format(domain, main_domain)


@pytest.fixture(scope='session')
def app_domain(device_domain):
    return 'platform.{0}'.format(device_domain)
