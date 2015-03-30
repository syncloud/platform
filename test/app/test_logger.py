import logging
from os import remove
import tempfile
import unittest
# from syncloud.app import logger
from syncloud.app import logger


class LoggingTest(unittest.TestCase):

    def setUp(self):
        self.test_log = tempfile.mktemp()

    def tearDown(self):
        remove(self.test_log)

    def test_log_duplicate_lines(self):
        logger.init(level=logging.DEBUG, console=False, filename=self.test_log)
        log = logger.get_logger('test')
        log.info('log1')
        log = logger.get_logger('test')
        log.info('log2')

        logList = open(self.test_log, 'r').read().splitlines()

        log1_count = len([line for line in logList if 'log1' in line])
        self.assertEqual(log1_count, 1)

        log2_count = len([line for line in logList if 'log2' in line])
        self.assertEqual(log2_count, 1)
