import logging
import os
from tempfile import mkstemp
from syncloud.app import logger
from syncloud.apache.setup import Setup
from syncloud.apache.system import System

logger.init(logging.DEBUG, True)


def test_generate_config():

    from_fh, from_file = mkstemp()
    with open(from_file, 'w') as f:
        f.write("test ${key}")

    to_fh, to_file = mkstemp()

    Setup(System()).generate_config(from_file, to_file, dict(key='123'))

    assert open(to_file, 'r').read() == 'test 123'
    assert os.path.exists("{}.conf".format(to_file))