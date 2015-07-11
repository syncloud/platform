import sys
import logging
from StringIO import StringIO
import pytest

from syncloud_app import runner


@pytest.fixture(scope="function")
def buffer():
    return StringIO()

@pytest.fixture(scope="function")
def logger(buffer):
    logger = logging.getLogger('logging_subprocess_test')
    logger.setLevel(logging.DEBUG)
    logHandler = logging.StreamHandler(buffer)
    formatter = logging.Formatter("%(levelname)s-%(message)s")
    logHandler.setFormatter(formatter)
    logger.addHandler(logHandler)
    return logger


def test_log_stdout(buffer, logger):
    runner.call([sys.executable, "-c", "print 'foo'"], logger)
    assert 'DEBUG-foo' in buffer.getvalue()


def test_log_stderr(buffer, logger):
    runner.call([sys.executable, "-c", 'import sys; sys.stderr.write("foo\\n")'], logger)
    assert 'ERROR-foo' in buffer.getvalue()


def test_custom_stdout_log_level(buffer, logger):
    runner.call([sys.executable, "-c", "print 'foo'"], logger, stdout_log_level=logging.INFO)
    assert 'INFO-foo' in buffer.getvalue()


def test_custom_stderr_log_level(buffer, logger):
    runner.call([sys.executable, "-c", 'import sys; sys.stderr.write("foo\\n")'],
                logger,
                stderr_log_level=logging.WARNING)
    assert 'WARNING-foo' in buffer.getvalue()