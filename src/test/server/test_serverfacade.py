import logging

from mock import MagicMock
import pytest

from syncloud_app import logger

logger.init(logging.DEBUG, True)


@pytest.fixture(scope="session")
def sam():
    sam = MagicMock()
    sam.update = MagicMock(return_value=True)
    sam.upgrade_all = MagicMock(return_value=True)
    return sam


@pytest.fixture(scope="session")
def insider():
    insider = MagicMock()
    insider.set_redirect_info = MagicMock(return_value=True)
    insider.acquire_domain = MagicMock(return_value=True)
    return insider

@pytest.fixture(scope="session")
def apache():
    apache = MagicMock()
    apache.activate = MagicMock()
    return apache