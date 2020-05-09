import glob
import os
from os.path import join
import tempfile
from subprocess import check_output

import ldap
from passlib import hash

from syncloudlib import gen, fs
from syncloudlib.logger import get_logger
import time

ldap_user_conf_dir = 'slapd.d'
DOMAIN="dc=syncloud,dc=org"


class LdapAuth:
    def __init__(self, platform_config, systemctl):
        self.systemctl = systemctl
        self.log = get_logger('ldap')
        self.config = platform_config
        self.user_conf_dir = join(self.config.data_dir(), ldap_user_conf_dir)
        self.ldap_root = '{0}/openldap'.format(self.config.app_dir())

    def installed(self):
        return os.path.isdir(join(self.config.data_dir(), ldap_user_conf_dir))

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

    def reset(self, name, user, password, email):

        self.systemctl.stop_service('platform.openldap')

        fs.removepath(self.user_conf_dir)

        files = glob.glob('{0}/openldap-data/*'.format(self.config.data_dir()))
        for f in files:
            os.remove(f)

        self.init()

        self.systemctl.start_service('platform.openldap')

        _, filename = tempfile.mkstemp()
        try:
            gen.transform_file('{0}/ldap/init.ldif'.format(self.config.config_dir()), filename, {
                'name': name,
                'user': user,
                'email': email,
                'password': make_secret(password)
            })

            self._init_db(filename)
        finally:
            os.remove(filename)

        check_output(generate_change_password_cmd(password), shell=True)

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


def generate_change_password_cmd(password):
    return 'echo "root:{0}" | chpasswd'.format(password.replace('"', '\\"').replace("$", "\\$"))


def to_ldap_dc(full_domain):
    return 'dc=' + ',dc='.join(full_domain.split('.'))


def authenticate(name, password):
    conn = ldap.initialize('ldap://localhost:389')
    try:
        conn.simple_bind_s('cn={0},ou=users,dc=syncloud,dc=org'.format(name), password)
    except Exception as e:
        conn.unbind()
        if 'desc' in e:
            raise Exception(e['desc'])
        else:
            raise Exception(str(e))


def make_secret(password):
    return hash.ldap_salted_sha1.hash(password)
