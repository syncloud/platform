import pytest
import requests

SYNCLOUD_INFO = 'syncloud.info'
DEVICE_USER = "user"
DEVICE_PASSWORD = "password"


def pytest_addoption(parser):
    parser.addoption("--email", action="store")
    parser.addoption("--password", action="store")
    parser.addoption("--domain", action="store")
    parser.addoption("--app-archive-path", action="store")
    parser.addoption("--release", action="store")
    parser.addoption("--installer", action="store")
    parser.addoption("--device-host", action="store")


@pytest.fixture(scope="session")
def auth(request):
    config = request.config
    return config.getoption("--email"), \
           config.getoption("--password"), \
           config.getoption("--domain"), \
           config.getoption("--release")


@pytest.fixture(scope="function")
def public_web_session(device_host):

    retry = 0
    retries = 5
    while retry < retries:
        try:
            session = requests.session()
            session.post('http://{0}/rest/login'.format(device_host), data={'name': DEVICE_USER, 'password': DEVICE_PASSWORD})
            assert session.get('http://{0}/rest/user'.format(device_host), allow_redirects=False).status_code == 200
            return session
        except Exception, e:
            retry += 1
            print(e.message)
            print('retry {0} of {1}'.format(retry, retries))


@pytest.fixture(scope='session')
def user_domain(main_domain):
    return 'platform.{0}'.format(main_domain)


@pytest.fixture(scope='session')
def main_domain(request):
    return '{0}.{1}'.format(request.config.getoption("--domain"), SYNCLOUD_INFO)

@pytest.fixture(scope='session')
def app_archive_path(request):
    return request.config.getoption("--app-archive-path")


@pytest.fixture(scope='session')
def installer(request):
    return request.config.getoption("--installer")


@pytest.fixture(scope='session')
def device_host(request):
    return request.config.getoption("--device-host")