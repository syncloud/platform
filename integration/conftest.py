from os.path import dirname, join, exists
import os
from syncloudlib.integration.conftest import *

DIR = dirname(__file__)


@pytest.fixture(scope="session")
def project_dir():
    return join(dirname(__file__), '..')


def pytest_addoption(parser):
    parser.addoption("--distro", action="store")
    parser.addoption("--domain", action="store")
    parser.addoption("--device-host", action="store")
    parser.addoption("--app-archive-path", action="store")
    parser.addoption("--app", action="store")
    parser.addoption("--ui-mode", action="store", default="desktop")
    parser.addoption("--device-user", action="store", default="user")
    parser.addoption("--build-number", action="store", default="local")
    parser.addoption("--browser", action="store", default="firefox")


@pytest.fixture(scope='session')
def distro(request):
    return request.config.getoption("--distro")


@pytest.fixture(scope="session")
def artifact_dir(project_dir, distro):
    dir = join(project_dir, 'artifact', distro)
    if not exists(dir):
        os.mkdir(dir)
    return dir
