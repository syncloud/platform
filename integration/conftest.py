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

@pytest.fixture(scope="function")
def public_web_session(device_host):

    retry = 0
    retries = 5
    while True:
        try:
            session = requests.session()
            session.post('https://{0}/rest/login'.format(device_host), verify=False, data={'name': DEVICE_USER, 'password': DEVICE_PASSWORD})
            assert session.get('https://{0}/rest/user'.format(device_host), verify=False, allow_redirects=False).status_code == 200
            return session
        except Exception, e:
            retry += 1
            if retry > retries:
                raise e
            print(e.message)
            print('retry {0} of {1}'.format(retry, retries))
    

