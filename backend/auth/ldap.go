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
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const ldapUserConfDir = "slapd.d"
const ldapUserDataDir = "openldap-data"
const Domain = "dc=syncloud,dc=org"
const UsersDn = "ou=users,dc=syncloud,dc=org"
const GroupsDn = "ou=groups,dc=syncloud,dc=org"
const AdminGroup = "syncloud"
const AdminGroupDn = "cn=syncloud,ou=groups,dc=syncloud,dc=org"
const posixIdStart = 2000

var emailRegexp = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var groupNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type User struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Admin    bool     `json:"admin"`
	Groups   []string `json:"groups"`
}

type Group struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type DomainProvider interface {
	GetDeviceDomain() string
}

type Service struct {
	snapService      SnapService
	runtimeConfigDir string
	userConfDir      string
	userDataDir      string
	ldapRoot         string
	configDir        string
	executor         cli.Executor
	passwordChanger  PasswordChanger
	domain           DomainProvider
	logger           *zap.Logger
}

type SnapService interface {
	Stop(name string) error
	Start(name string) error
}

func New(snapService SnapService, runtimeConfigDir string, appDir string, configDir string, executor cli.Executor, passwordChanger PasswordChanger, domain DomainProvider, logger *zap.Logger) *Service {

	return &Service{
		snapService:      snapService,
		runtimeConfigDir: runtimeConfigDir,
		userConfDir:      path.Join(runtimeConfigDir, ldapUserConfDir),
		userDataDir:      path.Join(runtimeConfigDir, ldapUserDataDir),
		ldapRoot:         path.Join(appDir, "openldap"),
		configDir:        configDir,
		executor:         executor,
		passwordChanger:  passwordChanger,
		domain:           domain,
		logger:           logger,
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

func (s *Service) ApplyConfig() error {
	if !s.Installed() {
		return nil
	}
	socket := path.Join(s.runtimeConfigDir, "openldap.socket")
	uri := fmt.Sprintf("ldapi://%s", strings.ReplaceAll(socket, "/", "%2F"))
	var err error
	for i := 0; i < 10; i++ {
		err = s.applyConfigOnce(uri)
		if err == nil {
			return nil
		}
		log.Printf("apply ldap config attempt %d failed: %s", i, err)
		time.Sleep(time.Second * 1)
	}
	return err
}

func (s *Service) applyConfigOnce(uri string) error {
	conn, err := ldap.DialURL(uri)
	if err != nil {
		return fmt.Errorf("ldapi connect: %w", err)
	}
	defer conn.Close()
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

func (s *Service) AuthenticateUser(username string, password string) error {
	conn, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		return fmt.Errorf("ldap connect: %w", err)
	}
	defer conn.Close()
	err = conn.Bind(fmt.Sprintf("cn=%s,ou=users,%s", username, Domain), password)
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	return nil
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

func (s *Service) IsAdmin(username string) (bool, error) {
	conn, err := s.rootBind()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	searchRequest := ldap.NewSearchRequest(
		AdminGroupDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf("(memberUid=%s)", ldap.EscapeFilter(username)),
		[]string{"memberUid"},
		nil)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false, fmt.Errorf("ldap search: %w", err)
	}
	return len(sr.Entries) > 0, nil
}

func (s *Service) rootBind() (*ldap.Conn, error) {
	conn, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		return nil, fmt.Errorf("ldap connect: %w", err)
	}
	err = conn.Bind(Domain, "syncloud")
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("ldap root bind: %w", err)
	}
	return conn, nil
}

const passwordMinLength = 8

func ValidatePassword(password string) error {
	if len(password) < passwordMinLength {
		return fmt.Errorf("password must be at least %d characters", passwordMinLength)
	}
	hasLetter := false
	hasDigit := false
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter {
		return fmt.Errorf("password must contain a letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain a number")
	}
	return nil
}

func (s *Service) ResolveEmail(username string, email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Sprintf("%s@%s", username, s.domain.GetDeviceDomain()), nil
	}
	if !emailRegexp.MatchString(email) {
		return "", fmt.Errorf("invalid email address: %s", email)
	}
	return email, nil
}

func userAttributes(username string, email string, id int) map[string][]string {
	idStr := strconv.Itoa(id)
	return map[string][]string{
		"objectClass":   {"person", "inetOrgPerson", "posixAccount", "simpleSecurityObject"},
		"cn":            {username},
		"sn":            {username},
		"givenName":     {username},
		"displayName":   {username},
		"uid":           {username},
		"uidNumber":     {idStr},
		"gidNumber":     {idStr},
		"homeDirectory": {"/home/" + username},
		"loginShell":    {"/bin/bash"},
		"mail":          {email},
	}
}

func (s *Service) AddUser(username string, password string, email string) error {
	if strings.TrimSpace(username) == "" {
		return fmt.Errorf("username is required")
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	resolvedEmail, err := s.ResolveEmail(username, email)
	if err != nil {
		return err
	}
	conn, err := s.rootBind()
	if err != nil {
		return err
	}
	defer conn.Close()

	id, err := s.nextUid(conn)
	if err != nil {
		return err
	}

	userDn := fmt.Sprintf("cn=%s,ou=users,%s", username, Domain)
	addReq := ldap.NewAddRequest(userDn, nil)
	for name, values := range userAttributes(username, resolvedEmail, id) {
		addReq.Attribute(name, values)
	}
	addReq.Attribute("userPassword", []string{makeSecret(password)})

	err = conn.Add(addReq)
	if err != nil {
		return fmt.Errorf("ldap add user: %w", err)
	}
	return nil
}

func (s *Service) SetUserEmail(username string, email string) error {
	resolvedEmail, err := s.ResolveEmail(username, email)
	if err != nil {
		return err
	}
	conn, err := s.rootBind()
	if err != nil {
		return err
	}
	defer conn.Close()

	userDn := fmt.Sprintf("cn=%s,ou=users,%s", username, Domain)
	modReq := ldap.NewModifyRequest(userDn, nil)
	modReq.Replace("mail", []string{resolvedEmail})
	if err := conn.Modify(modReq); err != nil {
		return fmt.Errorf("ldap set email: %w", err)
	}
	return nil
}

func (s *Service) SetPassword(username string, password string) error {
	if err := ValidatePassword(password); err != nil {
		return err
	}
	conn, err := s.rootBind()
	if err != nil {
		return err
	}
	defer conn.Close()

	userDn := fmt.Sprintf("cn=%s,ou=users,%s", username, Domain)
	modReq := ldap.NewModifyRequest(userDn, nil)
	modReq.Replace("userPassword", []string{makeSecret(password)})
	if err := conn.Modify(modReq); err != nil {
		return fmt.Errorf("ldap set password: %w", err)
	}
	return nil
}

func (s *Service) ListUsers() ([]User, error) {
	conn, err := s.rootBind()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	groups, err := s.listGroups(conn)
	if err != nil {
		return nil, err
	}
	membership := map[string][]string{}
	for _, group := range groups {
		for _, member := range group.Members {
			membership[member] = append(membership[member], group.Name)
		}
	}

	searchRequest := ldap.NewSearchRequest(
		UsersDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(objectClass=inetOrgPerson)",
		[]string{"cn", "mail"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("ldap list users: %w", err)
	}

	users := make([]User, 0, len(sr.Entries))
	for _, entry := range sr.Entries {
		username := entry.GetAttributeValue("cn")
		userGroups := membership[username]
		admin := false
		other := make([]string, 0, len(userGroups))
		for _, group := range userGroups {
			if group == AdminGroup {
				admin = true
			} else {
				other = append(other, group)
			}
		}
		users = append(users, User{
			Username: username,
			Email:    entry.GetAttributeValue("mail"),
			Admin:    admin,
			Groups:   other,
		})
	}
	return users, nil
}

func (s *Service) listGroups(conn *ldap.Conn) ([]Group, error) {
	searchRequest := ldap.NewSearchRequest(
		GroupsDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"cn", "memberUid", "gidNumber"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("ldap list groups: %w", err)
	}
	groups := make([]Group, 0, len(sr.Entries))
	for _, entry := range sr.Entries {
		groups = append(groups, Group{
			Name:    entry.GetAttributeValue("cn"),
			Members: entry.GetAttributeValues("memberUid"),
		})
	}
	return groups, nil
}

func (s *Service) ListGroups() ([]Group, error) {
	conn, err := s.rootBind()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return s.listGroups(conn)
}

func (s *Service) AddGroup(name string) error {
	if !groupNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid group name: %s", name)
	}
	conn, err := s.rootBind()
	if err != nil {
		return err
	}
	defer conn.Close()

	gid, err := s.nextGid(conn)
	if err != nil {
		return err
	}

	groupDn := fmt.Sprintf("cn=%s,%s", name, GroupsDn)
	addReq := ldap.NewAddRequest(groupDn, nil)
	addReq.Attribute("objectClass", []string{"posixGroup", "top"})
	addReq.Attribute("cn", []string{name})
	addReq.Attribute("gidNumber", []string{strconv.Itoa(gid)})
	if err := conn.Add(addReq); err != nil {
		return fmt.Errorf("ldap add group: %w", err)
	}
	return nil
}

func (s *Service) nextUid(conn *ldap.Conn) (int, error) {
	searchRequest := ldap.NewSearchRequest(
		UsersDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(objectClass=posixAccount)",
		[]string{"uidNumber"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return 0, fmt.Errorf("ldap uid scan: %w", err)
	}
	next := posixIdStart
	for _, entry := range sr.Entries {
		uid, err := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		if err == nil && uid >= next {
			next = uid + 1
		}
	}
	return next, nil
}

func (s *Service) nextGid(conn *ldap.Conn) (int, error) {
	searchRequest := ldap.NewSearchRequest(
		GroupsDn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"gidNumber"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return 0, fmt.Errorf("ldap gid scan: %w", err)
	}
	next := posixIdStart
	for _, entry := range sr.Entries {
		gid, err := strconv.Atoi(entry.GetAttributeValue("gidNumber"))
		if err == nil && gid >= next {
			next = gid + 1
		}
	}
	return next, nil
}

func (s *Service) RemoveGroup(name string) error {
	if name == AdminGroup {
		return fmt.Errorf("cannot remove admin group")
	}
	conn, err := s.rootBind()
	if err != nil {
		return err
	}
	defer conn.Close()

	groupDn := fmt.Sprintf("cn=%s,%s", name, GroupsDn)
	if err := conn.Del(ldap.NewDelRequest(groupDn, nil)); err != nil {
		return fmt.Errorf("ldap remove group: %w", err)
	}
	return nil
}

func (s *Service) SetAdmin(username string, admin bool) error {
	if !admin {
		conn, err := s.rootBind()
		if err != nil {
			return err
		}
		defer conn.Close()
		members, err := s.groupMembers(conn, AdminGroup)
		if err != nil {
			return err
		}
		if len(members) <= 1 && contains(members, username) {
			return fmt.Errorf("cannot remove the last admin")
		}
		return s.modifyMember(conn, AdminGroup, username, false)
	}
	return s.SetGroupMember(AdminGroup, username, true)
}

func (s *Service) SetGroupMember(group string, username string, member bool) error {
	conn, err := s.rootBind()
	if err != nil {
		return err
	}
	defer conn.Close()
	return s.modifyMember(conn, group, username, member)
}

func (s *Service) groupMembers(conn *ldap.Conn, group string) ([]string, error) {
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("cn=%s,%s", group, GroupsDn),
		ldap.ScopeBaseObject, ldap.DerefAlways, 0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"memberUid"},
		nil)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("ldap group members: %w", err)
	}
	if len(sr.Entries) < 1 {
		return nil, fmt.Errorf("group not found: %s", group)
	}
	return sr.Entries[0].GetAttributeValues("memberUid"), nil
}

func (s *Service) modifyMember(conn *ldap.Conn, group string, username string, member bool) error {
	members, err := s.groupMembers(conn, group)
	if err != nil {
		return err
	}
	present := contains(members, username)
	if member == present {
		return nil
	}
	groupDn := fmt.Sprintf("cn=%s,%s", group, GroupsDn)
	modReq := ldap.NewModifyRequest(groupDn, nil)
	if member {
		modReq.Add("memberUid", []string{username})
	} else {
		modReq.Delete("memberUid", []string{username})
	}
	if err := conn.Modify(modReq); err != nil {
		return fmt.Errorf("ldap modify group member: %w", err)
	}
	return nil
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func (s *Service) RemoveUser(username string) error {
	conn, err := s.rootBind()
	if err != nil {
		return err
	}
	defer conn.Close()

	userDn := fmt.Sprintf("cn=%s,ou=users,%s", username, Domain)
	delReq := ldap.NewDelRequest(userDn, nil)
	err = conn.Del(delReq)
	if err != nil {
		return fmt.Errorf("ldap delete user: %w", err)
	}
	return nil
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
