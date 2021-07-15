import syncloudlib.integration.conftest
from os.path import dirname, join, exists
import os
from syncloudlib.integration.conftest import *

DIR = dirname(__file__)


@pytest.fixture(scope="session")
def project_dir():
    return join(dirname(__file__), '..')


def pytest_addoption(parser):
    syncloudlib.integration.conftest.pytest_addoption(parser)
    parser.addoption("--distro", action="store")
    parser.addoption("--arch", action="store", default="unset-arch")


@pytest.fixture(scope='session')
def distro(request):
    return request.config.getoption("--distro")


@pytest.fixture(scope='session')
def arch(request):
    return request.config.getoption("--arch")


@pytest.fixture(scope='session')
def full_domain(domain, main_domain):
    return '{}.{}'.format(domain, main_domain)


@pytest.fixture(scope="session")
def artifact_dir(project_dir, distro):
    dir = join(project_dir, 'artifact', distro)
    if not exists(dir):
        os.mkdir(dir)
    return dir
