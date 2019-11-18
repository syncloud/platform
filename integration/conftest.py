from os.path import dirname, join
from syncloudlib.integration.conftest import *

DIR = dirname(__file__)


@pytest.fixture(scope="session")
def project_dir():
    return join(dirname(__file__), '..')
