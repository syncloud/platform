import glob
import hashlib
import json
import os
from os.path import join
import tempfile
from subprocess import check_output

import ldap

from syncloud_app import util
from syncloud_app.logger import get_logger
from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.systemd.systemctl import stop_service, start_service
from syncloud_platform.tools import app

ldap_user_conf_dir='slapd.d'
platform_user = 'platform'


class Auth:
    def __init__(self):
        self.logger = get_logger('ldap')
        self.config = PlatformConfig()

    def installed(self):
        data_dir = app.get_app_data_root('platform', platform_user)
        return os.path.isdir(join(data_dir, ldap_user_conf_dir))

    def reset(self, user, password):

        data_dir = app.get_app_data_root('platform', platform_user)
        user_conf_dir = app.create_data_dir(data_dir, ldap_user_conf_dir, platform_user, remove_existing=True)

        stop_service('platform-openldap')

        files = glob.glob('/opt/app/platform/openldap/var/openldap-data/*')
        for f in files:
            os.remove(f)

        init_script = '{0}/ldap/slapd.ldif'.format(self.config.config_dir())
        ldap_root = '{0}/openldap'.format(self.config.app_dir())

        check_output(
            '{0}/sbin/slapadd -F {1} -b "cn=config" -l {2}'.format(ldap_root, user_conf_dir, init_script), shell=True)

        check_output('chown -R {0}. {1}'.format(platform_user, user_conf_dir), shell=True)

        start_service('platform-openldap')

        fd, filename = tempfile.mkstemp()
        util.transform_file('{0}/ldap/init.ldif'.format(self.config.config_dir()), filename, {
            'user': user,
            'password': make_secret(password)
        })

        check_output('{0}/bin/ldapadd -Y EXTERNAL -H ldapi:/// -f {1}'.format(ldap_root, filename), shell=True)

def to_ldap_dc(full_domain):
    return 'dc=' + ',dc='.join(full_domain.split('.'))


def authenticate(name, password):
    conn = ldap.initialize('ldap://localhost:389')
    try:
        conn.simple_bind_s('cn={0},ou=users,dc=syncloud,dc=org'.format(name), password)
    except Exception, e:
        conn.unbind()
        if 'desc' in e.message:
            raise Exception(e.message['desc'])
        else:
            raise Exception(e.message)

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