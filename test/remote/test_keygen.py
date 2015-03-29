import logging
import tempfile
from os import path

from syncloud.app import logger

from syncloud.remote.keygen import KeyGen

logger.init()

class TestRemote():

    def test_key_gen(self):

        logger.console = True
        logger.level = logging.DEBUG
        private, public = KeyGen().generate('rsa')

        assert len(private) > 0
        assert len(public) > 0

    def test_key_gen_with_bits(self):

        logger.console = True
        logger.level = logging.DEBUG
        private, public = KeyGen().generate('rsa', 2048)

        assert len(private) > 0
        assert len(public) > 0

    def test_key_gen_overwrite(self):

        tempdir = tempfile.mkdtemp()
        tempkey = path.join(tempdir, 'key')

        key_gen = KeyGen()
        key_gen.generate_into_file('rsa', tempkey, 2048, True)

        private1 = key_gen.read(tempkey)
        public1 = key_gen.read(tempkey + '.pub')

        key_gen.generate_into_file('rsa', tempkey, 2048, True)

        private2 = key_gen.read(tempkey)
        public2 = key_gen.read(tempkey + '.pub')

        assert private1 != private2
        assert public1 != public2
