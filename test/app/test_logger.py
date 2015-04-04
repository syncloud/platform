import logging
from os import remove
import tempfile
import pytest
from syncloud.app import logger


@pytest.fixture(scope="module")
def testfile(request):
    tfile = tempfile.mktemp()
    def fin():
        print ("teardown tfile")
        remove(tfile)
    request.addfinalizer(fin)
    return tfile


def test_log_duplicate_lines(testfile):
    logger.init(level=logging.DEBUG, console=False, filename=testfile)
    log = logger.get_logger('test')
    log.info('log1')
    log = logger.get_logger('test')
    log.info('log2')

    logList = open(testfile, 'r').read().splitlines()

    log1_count = len([line for line in logList if 'log1' in line])
    assert log1_count == 1

    log2_count = len([line for line in logList if 'log2' in line])
    assert log2_count == 1
