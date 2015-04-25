import logging
import logging.handlers
import sys

factory_instance = None


class WhitespaceRemovingFormatter(logging.Formatter):
    def format(self, record):
        record.msg = clean(record.msg)
        return super(WhitespaceRemovingFormatter, self).format(record)


def clean(message):
        return str(message).strip()


class LoggerFactory:
    def __init__(self, level, console, filename, line_format):
        self.level = level
        self.console = console
        self.filename = filename
        self.formatter = WhitespaceRemovingFormatter(line_format)
        if filename:
            self.file_handler = logging.handlers.RotatingFileHandler(self.filename, maxBytes=1000000, backupCount=5)
            self.file_handler.setFormatter(self.formatter)
        if console:
            self.console_handler = logging.StreamHandler(sys.stdout)
            self.console_handler.setFormatter(self.formatter)

    def get_logger(self, name):
        logger = logging.getLogger(name)
        logger.setLevel(self.level)
        if self.filename:
            if self.file_handler not in logger.handlers:
                logger.addHandler(self.file_handler)
        if self.console:
            if self.console_handler not in logger.handlers:
                logger.addHandler(self.console_handler)
        return logger

default_format = '%(asctime)s - %(name)s - %(levelname)s - %(message)s'


def init(level=logging.INFO, console=False, filename=None, line_format=default_format):
    global factory_instance
    factory_instance = LoggerFactory(level, console, filename, line_format)


def get_logger(name):
    global factory_instance
    if not factory_instance:
        raise Exception('Logging is not initialized')
    return factory_instance.get_logger(name)
