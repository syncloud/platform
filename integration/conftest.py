import syncloudlib.integration.conftest
from os.path import dirname, join, exists
import os
from syncloudlib.integration.conftest import *

DIR = dirname(__file__)


@pytest.fixture(scope="session")
def project_dir():
    return join(dirname(__file__), '..')


@pytest.fixture(scope='session')
def main_domain():
    return 'syncloud.info'


@pytest.fixture(scope='session')
def full_domain(domain, main_domain):
    return '{}.{}'.format(domain, main_domain)
