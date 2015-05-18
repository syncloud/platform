import glob
import hashlib
import os
from os.path import join
import tempfile
from syncloud.app import util
from syncloud.app.logger import get_logger
from syncloud.tools.facade import Facade
from syncloud.tools.service import Service
from syncloud.app import runner


class Ldap():
    def __init__(self):
        self.logger = get_logger('ldap')
        self.service = Service()
        tools_facade = Facade()
        config_dir = join(tools_facade.usr_local_dir(), 'syncloud', 'ldap', 'config')
        self.rootdn_ldif = join(config_dir, 'rootdn.ldif')
        self.init_ldif = join(config_dir, 'init.ldif')

    def reset(self, full_domain, user, password):

        self.service.stop('slapd')

        files = glob.glob('/var/lib/ldap/*')
        for f in files:
            os.remove(f)

        self.service.start('slapd')

        dn = to_ldap_dc(full_domain)

        fd, filename = tempfile.mkstemp()
        secret = make_secret(password)
        util.transform_file(self.rootdn_ldif, filename, {
            'dn': dn,
            'password': secret
        })
        runner.call('ldapmodify -Y EXTERNAL -H ldapi:/// -f {0}'.format(filename), self.logger, shell=True)

        fd, filename = tempfile.mkstemp()
        util.transform_file(self.init_ldif, filename, {
            'dn': dn,
            'user': user,
            'password': secret
        })
        runner.call('ldapadd -Y EXTERNAL -H ldapi:/// -f {0}'.format(filename), self.logger, shell=True)


def to_ldap_dc(full_domain):
    return 'dc=' + ',dc='.join(full_domain.split('.'))


#https://gist.github.com/rca/7217540
def make_secret(password):
    """
    Encodes the given password as a base64 SSHA hash+salt buffer
    """
    salt = os.urandom(4)

    # hash the password and append the salt
    sha = hashlib.sha1(password)
    sha.update(salt)

    # create a base64 encoded string of the concatenated digest + salt
    digest_salt_b64 = '{}{}'.format(sha.digest(), salt).encode('base64').strip()

    # now tag the digest above with the {SSHA} tag
    tagged_digest_salt = '{{SSHA}}{}'.format(digest_salt_b64)

    return tagged_digest_salt