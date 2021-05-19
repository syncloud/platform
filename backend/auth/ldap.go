package auth

import (
	"fmt"
	"strings"
)

const ldapUserConfDir = "slapd.d"
const Domain = "dc=syncloud,dc=org"

type LdapAuth struct {
}

func New() *LdapAuth {
	//self.systemctl = systemctl
	//self.log = get_logger('ldap')
	//self.config = platform_config
	//self.user_conf_dir = join(self.config.data_dir(), ldap_user_conf_dir)
	//self.ldap_root = '{0}/openldap'.format(self.config.app_dir())
	return &LdapAuth{}
}

func (ldap *LdapAuth) Installed() bool {
	//return os.path.isdir(join(self.config.data_dir(), ldap_user_conf_dir))
	return false
}

func (ldap *LdapAuth) Init() {
	/*
	    if self.installed():
	       self.log.info('ldap config already initialized')
	       return

	   self.log.info('initializing ldap config')
	   fs.makepath(self.user_conf_dir)
	   init_script = '{0}/ldap/slapd.ldif'.format(self.config.config_dir())

	   check_output(
	       '{0}/sbin/slapadd.sh -F {1} -b "cn=config" -l {2}'.format(
	           self.ldap_root, self.user_conf_dir, init_script), shell=True)

	*/
}
func (ldap *LdapAuth) Reset(name string, user string, password string, email string) {
	/*
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

	*/
}

func (ldap *LdapAuth) initDb(filename string) {
	/*        success = False
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
	*/
}

func (ldap *LdapAuth) ldapAdd(filename string, bindDn string) {
	/*        bind_dn_option = ''
	if bind_dn:
	    bind_dn_option = '-D "{0}"'.format(bind_dn)
	check_output('{0}/bin/ldapadd.sh -x -w syncloud {1} -f {2}'.format(
	            self.ldap_root, bind_dn_option, filename), shell=True)
	*/
}

func GenerateChangePasswordCmd(password string) string {
	//TODO: fix me, we should not depend on bash limitations for passwords
	//return fmt.Sprintf("echo \"root:%s\" | chpasswd", strings.ReplaceAll()password.replace('"', '\\"').replace("$", "\\$"))
	return ""
}

func ToLdapDc(fullDomain string) string {
	return fmt.Sprintf("dc=%s", strings.Join(strings.Split(fullDomain, "."), ",dc="))
}

func Authenticate(name string, password string) {
	/*    conn = ldap.initialize('ldap://localhost:389')
	try:
	    conn.simple_bind_s('cn={0},ou=users,dc=syncloud,dc=org'.format(name), password)
	except ldap.INVALID_CREDENTIALS:
	    conn.unbind()
	    raise Exception('Invalid credentials')
	except Exception as e:
	    conn.unbind()
	    raise Exception(str(e))

	*/
}

func makeSecret(password string) string {
	//return hash.ldap_salted_sha1.hash(password)
	return ""
}
