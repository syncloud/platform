import sys
import unittest
import logging
from StringIO import StringIO

from syncloud.app import runner

class LoggingSubprocessTest(unittest.TestCase):

    def setUp(self):
        self.buffer = StringIO()
        self.logger = logging.getLogger('logging_subprocess_test')
        self.logger.setLevel(logging.DEBUG)
        self.logHandler = logging.StreamHandler(self.buffer)
        formatter = logging.Formatter("%(levelname)s-%(message)s")
        self.logHandler.setFormatter(formatter)
        self.logger.addHandler(self.logHandler)

    def test_log_stdout(self):
        runner.call([sys.executable, "-c", "print 'foo'"], self.logger)
        self.assertIn('DEBUG-foo', self.buffer.getvalue())

    def test_log_stderr(self):
        runner.call([sys.executable, "-c", 'import sys; sys.stderr.write("foo\\n")'], self.logger)
        self.assertIn('ERROR-foo', self.buffer.getvalue())

    def test_custom_stdout_log_level(self):
        runner.call([sys.executable, "-c", "print 'foo'"], self.logger, stdout_log_level=logging.INFO)
        self.assertIn('INFO-foo', self.buffer.getvalue())

    def test_custom_stderr_log_level(self):
        runner.call([sys.executable, "-c", 'import sys; sys.stderr.write("foo\\n")'],
                    self.logger,
                    stderr_log_level=logging.WARNING)
        self.assertIn('WARNING-foo', self.buffer.getvalue())