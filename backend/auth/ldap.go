package auth

import (
	"fmt"
	"github.com/syncloud/platform/snap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const ldapUserConfDir = "slapd.d"
const ldapUserDataDir = "openldap-data"
const Domain = "dc=syncloud,dc=org"

type Service struct {
	snapService *snap.Service
	userConfDir string
	userDataDir string
	ldapRoot    string
	configDir   string
}

func New(snapService *snap.Service, dataDir string, appDir string, configDir string) *Service {

	return &Service{
		snapService: snapService,
		userConfDir: path.Join(dataDir, ldapUserConfDir),
		userDataDir: path.Join(dataDir, ldapUserDataDir),
		ldapRoot:    path.Join(appDir, "openldap"),
		configDir:   configDir,
	}
}

func (l *Service) Installed() bool {
	_, err := os.Stat(l.userConfDir)
	return err == nil
}

func (l *Service) Init() error {
	if l.Installed() {
		log.Println("ldap config already initialized")
		return nil
	}
	log.Println("initializing ldap config")
	err := os.MkdirAll(l.userConfDir, 755)
	if err != nil {
		return err
	}

	initScript := path.Join(l.configDir, "ldap", "slapd.ldif")

	cmd := path.Join(l.ldapRoot, "sbin", "slapadd.sh")
	output, err := exec.Command(cmd, "-F", l.userConfDir, "-b", "cn=config", "-l", initScript).CombinedOutput()
	if err != nil {
		return err
	}
	log.Println(output)
	return nil
}

func (l *Service) Reset(name string, user string, password string, email string) error {
	log.Println("resetting ldap")

	err := l.snapService.Stop("platform.openldap")
	if err != nil {
		return err
	}
	err = os.RemoveAll(l.userConfDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(l.userDataDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(l.userDataDir, 755)
	if err != nil {
		return err
	}

	err = l.Init()
	if err != nil {
		return err
	}
	err = l.snapService.Start("platform.openldap")
	if err != nil {
		return err
	}

	passwordHash, err := makeSecret(password)
	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	file, err := ioutil.ReadFile(path.Join(l.configDir, "ldap", "init.ldif"))
	if err != nil {
		return err
	}
	ldif := string(file)
	ldif = strings.ReplaceAll(ldif, "${name}", name)
	ldif = strings.ReplaceAll(ldif, "${user}", user)
	ldif = strings.ReplaceAll(ldif, "${email}", email)
	ldif = strings.ReplaceAll(ldif, "${password}", *passwordHash)
	err = ioutil.WriteFile(tmpFile.Name(), []byte(ldif), 644)
	if err != nil {
		return err
	}

	err = l.initDb(tmpFile.Name())
	if err != nil {
		return err
	}

	err = ChangeSystemPassword(password)
	return err
}

func (l *Service) initDb(filename string) error {
	return l.ldapAdd(filename, Domain)
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

func (l *Service) ldapAdd(filename string, bindDn string) error {
	cmd := path.Join(l.ldapRoot, "bin", "ldapadd.sh")
	_, err := exec.Command(cmd, "-x", "-w", "syncloud", "-D", bindDn, "-f", filename).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func ChangeSystemPassword(password string) error {
	cmd := exec.Command("chpasswd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	_, err = io.WriteString(stdin, fmt.Sprintf("root:%s", password))
	return err
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

func makeSecret(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashString := string(hash)
	return &hashString, nil
}
