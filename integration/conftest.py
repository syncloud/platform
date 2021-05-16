from os.path import dirname, join
from syncloudlib.integration.conftest import *

DIR = dirname(__file__)


@pytest.fixture(scope="session")
def project_dir():
    return join(dirname(__file__), '..')


@pytest.fixture(scope="session")
def redirect_api_url(main_domain):
    return 'https://api.{}'.format(main_domain)
