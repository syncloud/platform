import glob
import hashlib
import os
from os.path import join
import tempfile
from syncloud.app import util
from syncloud.app.logger import get_logger
from syncloud.systemd.systemctl import stop_service, start_service
from syncloud.tools.facade import Facade
from syncloud.tools.service import Service
from syncloud.app import runner
import ldap


class Auth:
    def __init__(self):
        self.logger = get_logger('ldap')
        self.service = Service()
        tools_facade = Facade()
        config_dir = join(tools_facade.usr_local_dir(), 'syncloud', 'ldap', 'config')
        self.rootdn_ldif = join(config_dir, 'rootdn.ldif')
        self.init_ldif = join(config_dir, 'init.ldif')

    def reset(self, full_domain, user, password):

        stop_service('platform-openldap')

        files = glob.glob('/opt/app/platform/openldap/var/openldap-data/*')
        for f in files:
            os.remove(f)

        start_service('platform-openldap')

        dn = to_ldap_dc(full_domain)

        fd, filename = tempfile.mkstemp()
        secret = make_secret(password)
        util.transform_file(self.rootdn_ldif, filename, {
            'dn': dn,
            'password': secret
        })
        exit_code = runner.call('/opt/app/platform/openldap/bin/ldapmodify -Y EXTERNAL -H ldapi:/// -f {0}'.format(filename), self.logger, shell=True)
        if not exit_code == 0:
            raise Exception("Non zero exit code: {0}".format(exit_code))

        fd, filename = tempfile.mkstemp()
        util.transform_file(self.init_ldif, filename, {
            'dn': dn,
            'user': user,
            'password': secret
        })
        exit_code = runner.call('/opt/app/platform/openldap/bin/ldapadd -Y EXTERNAL -H ldapi:/// -f {0}'.format(filename), self.logger, shell=True)
        if not exit_code == 0:
            raise Exception("Non zero exit code: {0}".format(exit_code))

def to_ldap_dc(full_domain):
    return 'dc=' + ',dc='.join(full_domain.split('.'))


def authenticate(name, password, full_domain_name):
    conn = ldap.initialize('ldap://localhost:389')
    try:
        conn.simple_bind_s('cn={0},ou=users,{1}'.format(name, to_ldap_dc(full_domain_name)), password)
    except Exception, e:
        conn.unbind()
        raise e


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
    digest_salt_b64 = '{0}{1}'.format(sha.digest(), salt).encode('base64').strip()

    # now tag the digest above with the {SSHA} tag
    tagged_digest_salt = '{{SSHA}}{0}'.format(digest_salt_b64)

    return tagged_digest_salt