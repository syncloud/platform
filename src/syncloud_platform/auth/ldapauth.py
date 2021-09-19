import ldap
import os
import time
from passlib import hash
from subprocess import check_output
from syncloudlib import fs
from syncloudlib.logger import get_logger

ldap_user_conf_dir = '/var/snap/platform/current/slapd.d'
DOMAIN = "dc=syncloud,dc=org"


class LdapAuth:
    def __init__(self, platform_config, systemctl):
        self.systemctl = systemctl
        self.log = get_logger('ldap')
        self.config = platform_config
        self.user_conf_dir = ldap_user_conf_dir
        self.ldap_root = '{0}/openldap'.format(self.config.app_dir())

    def installed(self):
        return os.path.isdir(self.user_conf_dir)

    def init(self):
        if self.installed():
            self.log.info('ldap config already initialized')
            return

        self.log.info('initializing ldap config')
        fs.makepath(self.user_conf_dir)
        init_script = '{0}/ldap/slapd.ldif'.format(self.config.config_dir())

        check_output(
            '{0}/sbin/slapadd.sh -F {1} -b "cn=config" -l {2}'.format(
                self.ldap_root, self.user_conf_dir, init_script), shell=True)

    def _init_db(self, filename):
        success = False
        for i in range(0, 3):
            try:
                self.ldapadd(filename, DOMAIN)
                success = True
                break
            except Exception as e:
                self.log.warn(str(e))
                self.log.warn("probably ldap is still starting, will retry {0}".format(i))
                time.sleep(1)

        if not success:
            raise Exception("Unable to initialize ldap db")

    def ldapadd(self, filename, bind_dn=None):
        bind_dn_option = ''
        if bind_dn:
            bind_dn_option = '-D "{0}"'.format(bind_dn)
        check_output('{0}/bin/ldapadd.sh -x -w syncloud {1} -f {2}'.format(
            self.ldap_root, bind_dn_option, filename), shell=True)

    def authenticate(self, name, password):
        conn = ldap.initialize('ldap://localhost:389')
        try:
            conn.simple_bind_s('cn={0},ou=users,dc=syncloud,dc=org'.format(name), password)
            self.log.info('{0} authenticated'.format(name))
        except ldap.INVALID_CREDENTIALS:
            self.log.warn('{0} not authenticated'.format(name))
            conn.unbind()
            raise Exception('Invalid credentials')
        except Exception as e:
            self.log.warn('{0} not authenticated'.format(name))
            conn.unbind()
            raise Exception(str(e))


def make_secret(password):
    return hash.ldap_salted_sha1.hash(password)
