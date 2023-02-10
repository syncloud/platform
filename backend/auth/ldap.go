package auth

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/syncloud/platform/cli"
	"go.uber.org/zap"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const ldapUserConfDir = "slapd.d"
const ldapUserDataDir = "openldap-data"
const Domain = "dc=syncloud,dc=org"
const AdminGroupDn = "cn=syncloud,ou=groups,dc=syncloud,dc=org"

type Service struct {
	snapService     SnapService
	userConfDir     string
	userDataDir     string
	ldapRoot        string
	configDir       string
	executor        cli.Executor
	passwordChanger PasswordChanger
	logger          *zap.Logger
}

type SnapService interface {
	Stop(name string) error
	Start(name string) error
}

func New(snapService SnapService, runtimeConfigDir string, appDir string, configDir string, executor cli.Executor, passwordChanger PasswordChanger, logger *zap.Logger) *Service {

	return &Service{
		snapService:     snapService,
		userConfDir:     path.Join(runtimeConfigDir, ldapUserConfDir),
		userDataDir:     path.Join(runtimeConfigDir, ldapUserDataDir),
		ldapRoot:        path.Join(appDir, "openldap"),
		configDir:       configDir,
		executor:        executor,
		passwordChanger: passwordChanger,
		logger:          logger,
	}
}

func (s *Service) Installed() bool {
	_, err := os.Stat(s.userConfDir)
	return err == nil
}

func (s *Service) Init() error {
	if s.Installed() {
		log.Println("ldap config already initialized")
		return nil
	}
	log.Println("initializing ldap config")
	err := os.MkdirAll(s.userConfDir, 755)
	if err != nil {
		return err
	}

	initScript := path.Join(s.configDir, "ldap", "slapd.ldif")

	cmd := path.Join(s.ldapRoot, "sbin", "slapadd.sh")
	output, err := s.executor.CombinedOutput(cmd, "-F", s.userConfDir, "-b", "cn=config", "-l", initScript)
	if err != nil {
		return err
	}
	log.Println(string(output))
	return nil
}

func (s *Service) Reset(name string, user string, password string, email string) error {
	log.Println("resetting ldap")

	err := s.snapService.Stop("platform.openldap")
	if err != nil {
		return err
	}
	err = os.RemoveAll(s.userConfDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(s.userDataDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(s.userDataDir, 755)
	if err != nil {
		return err
	}

	err = s.Init()
	if err != nil {
		return err
	}
	err = s.snapService.Start("platform.openldap")
	if err != nil {
		return err
	}

	passwordHash := makeSecret(password)

	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	file, err := os.ReadFile(path.Join(s.configDir, "ldap", "init.ldif"))
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

	err = s.initDb(tmpFile.Name())
	if err != nil {
		return err
	}

	err = s.passwordChanger.Change(password)
	return err
}

func (s *Service) initDb(filename string) error {
	var err error
	for i := 0; i < 10; i++ {
		err = s.ldapAdd(filename, Domain)
		if err == nil {
			break
		}
		log.Println(err)
		log.Printf("probably ldap is still starting, will retry %d\n", i)
		time.Sleep(time.Second * 1)
	}
	return err
}

func (s *Service) ldapAdd(filename string, bindDn string) error {
	cmd := path.Join(s.ldapRoot, "bin", "ldapadd.sh")
	output, err := s.executor.CombinedOutput(cmd, "-x", "-w", "syncloud", "-D", bindDn, "-f", filename)
	log.Printf("ldapadd output: %s", output)
	return err
}

func (s *Service) Authenticate(username string, password string) (bool, error) {
	conn, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		return false, err
	}
	defer conn.Close()
	err = conn.Bind(fmt.Sprintf("cn=%s,ou=users,dc=syncloud,dc=org", username), password)
	if err != nil {
		s.logger.Error("ldap error", zap.Error(err))
		return false, err
	}

	searchRequest := ldap.NewSearchRequest(
		AdminGroupDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf("(memberUid=%s)", username),
		[]string{"memberUid"},
		nil)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false, err
	}

	if len(sr.Entries) < 1 {
		return false, fmt.Errorf("not admin (must be part of syncloud group)")
	}
	return true, nil
}

func makeSecret(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	salt := make([]byte, 4)
	_, err := rand.Read(salt)
	if err != nil {
		log.Printf("unable to generate password salt: %s", err)
		salt = []byte("salt")
	}
	hasher.Write(salt)
	hash := hasher.Sum(nil)
	hashWithSalt := append(hash, salt...)
	encodedHash := base64.StdEncoding.EncodeToString(hashWithSalt)
	return fmt.Sprintf("{SSHA}%s", encodedHash)
}
