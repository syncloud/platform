from os.path import dirname

from syncloudlib.integration.conftest import *

DIR = dirname(__file__)


@pytest.fixture(scope="session")
def project_dir():
    return join(dirname(__file__), '..')


@pytest.fixture(scope='session')
def main_domain():
    return 'redirect'


@pytest.fixture(scope='session')
def full_domain(domain, main_domain):
    return '{}.{}'.format(domain, main_domain)
