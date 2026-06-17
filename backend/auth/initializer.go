package auth

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/syncloud/platform/cli"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const ldapUserConfDir = "slapd.d"
const ldapUserDataDir = "openldap-data"
const Domain = "dc=syncloud,dc=org"
const UsersDn = "ou=users,dc=syncloud,dc=org"
const GroupsDn = "ou=groups,dc=syncloud,dc=org"
const AdminGroup = "syncloud"
const AdminGroupDn = "cn=syncloud,ou=groups,dc=syncloud,dc=org"
const posixIdStart = 2000

type Initializer struct {
	snapService      SnapService
	runtimeConfigDir string
	userConfDir      string
	userDataDir      string
	ldapRoot         string
	configDir        string
	executor         cli.Executor
	ldapClient       *LdapClient
	passwordChanger  PasswordChanger
	passwordHasher   *PasswordHasher
}

type SnapService interface {
	Stop(name string) error
	Start(name string) error
}

func NewInitializer(snapService SnapService, runtimeConfigDir string, appDir string, configDir string, executor cli.Executor, ldapClient *LdapClient, passwordChanger PasswordChanger, passwordHasher *PasswordHasher) *Initializer {

	return &Initializer{
		snapService:      snapService,
		runtimeConfigDir: runtimeConfigDir,
		userConfDir:      path.Join(runtimeConfigDir, ldapUserConfDir),
		userDataDir:      path.Join(runtimeConfigDir, ldapUserDataDir),
		ldapRoot:         path.Join(appDir, "openldap"),
		configDir:        configDir,
		executor:         executor,
		ldapClient:       ldapClient,
		passwordChanger:  passwordChanger,
		passwordHasher:   passwordHasher,
	}
}

func (i *Initializer) Installed() bool {
	_, err := os.Stat(i.userConfDir)
	return err == nil
}

func (i *Initializer) Init() error {
	if i.Installed() {
		log.Println("ldap config already initialized")
		return nil
	}
	log.Println("initializing ldap config")
	err := os.MkdirAll(i.userConfDir, 755)
	if err != nil {
		return err
	}

	initScript := path.Join(i.configDir, "ldap", "slapd.ldif")

	cmd := path.Join(i.ldapRoot, "sbin", "slapadd.sh")
	output, err := i.executor.CombinedOutput(cmd, "-F", i.userConfDir, "-b", "cn=config", "-l", initScript)
	if err != nil {
		return err
	}
	log.Println(string(output))
	return nil
}

func (i *Initializer) ApplyConfig() error {
	if !i.Installed() {
		return nil
	}
	socket := path.Join(i.runtimeConfigDir, "openldap.socket")
	uri := fmt.Sprintf("ldapi://%s", strings.ReplaceAll(socket, "/", "%2F"))
	var err error
	for attempt := 0; attempt < 10; attempt++ {
		err = i.applyConfigOnce(uri)
		if err == nil {
			return nil
		}
		log.Printf("apply ldap config attempt %d failed: %s", attempt, err)
		time.Sleep(time.Second * 1)
	}
	return err
}

func (i *Initializer) applyConfigOnce(uri string) error {
	conn, err := ldap.DialURL(uri)
	if err != nil {
		return fmt.Errorf("ldapi connect: %w", err)
	}
	defer i.ldapClient.Disconnect(conn)
	if err := conn.ExternalBind(); err != nil {
		return fmt.Errorf("ldapi external bind: %w", err)
	}
	req := ldap.NewModifyRequest("cn=config", nil)
	req.Replace("olcLogLevel", []string{"none"})
	if err := conn.Modify(req); err != nil {
		return fmt.Errorf("ldap modify cn=config: %w", err)
	}
	return nil
}

func (i *Initializer) Reset(name string, user string, password string, email string) error {
	log.Println("resetting ldap")

	err := i.snapService.Stop("platform.openldap")
	if err != nil {
		return err
	}
	err = os.RemoveAll(i.userConfDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(i.userDataDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(i.userDataDir, 755)
	if err != nil {
		return err
	}

	err = i.Init()
	if err != nil {
		return err
	}
	err = i.snapService.Start("platform.openldap")
	if err != nil {
		return err
	}

	passwordHash := i.passwordHasher.Hash(password)

	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	file, err := os.ReadFile(path.Join(i.configDir, "ldap", "init.ldif"))
	if err != nil {
		return err
	}
	ldif := string(file)
	ldif = strings.ReplaceAll(ldif, "${name}", name)
	ldif = strings.ReplaceAll(ldif, "${user}", user)
	ldif = strings.ReplaceAll(ldif, "${email}", email)
	ldif = strings.ReplaceAll(ldif, "${password}", passwordHash)
	err = os.WriteFile(tmpFile.Name(), []byte(ldif), 644)
	if err != nil {
		return err
	}

	err = i.initDb(tmpFile.Name())
	if err != nil {
		return err
	}

	err = i.passwordChanger.Change(password)
	return err
}

func (i *Initializer) initDb(filename string) error {
	var err error
	for attempt := 0; attempt < 10; attempt++ {
		err = i.ldapAdd(filename, Domain)
		if err == nil {
			break
		}
		log.Println(err)
		log.Printf("probably ldap is still starting, will retry %d\n", attempt)
		time.Sleep(time.Second * 1)
	}
	return err
}

func (i *Initializer) ldapAdd(filename string, bindDn string) error {
	cmd := path.Join(i.ldapRoot, "bin", "ldapadd.sh")
	output, err := i.executor.CombinedOutput(cmd, "-x", "-w", "syncloud", "-D", bindDn, "-f", filename)
	log.Printf("ldapadd output: %s", output)
	return err
}
