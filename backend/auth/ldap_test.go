package auth

import (
	"fmt"
	"github.com/syncloud/platform/log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SnapServiceStub struct {
}

func (s SnapServiceStub) Stop(_ string) error {
	return nil
}

func (s SnapServiceStub) Start(_ string) error {
	return nil
}

type ExecutorStub struct {
	executions []string
}

func (e *ExecutorStub) CombinedOutput(name string, arg ...string) ([]byte, error) {
	e.executions = append(e.executions, fmt.Sprintf("%s %s", name, strings.Join(arg, " ")))
	return []byte(""), nil
}

type PasswordChangerStub struct {
	changed bool
}

func (p *PasswordChangerStub) Change(_ string) error {
	p.changed = true
	return nil
}

type DomainProviderStub struct {
	domain string
}

func (d DomainProviderStub) GetDeviceDomain() string {
	return d.domain
}

func TestMakeSecret(t *testing.T) {
	secret := makeSecret("syncloud")
	assert.Greater(t, len(secret), 1)
}

func newTestService(domain string) *Service {
	return New(&SnapServiceStub{}, t1TempDir(), t1TempDir(), t1TempDir(), &ExecutorStub{}, &PasswordChangerStub{}, DomainProviderStub{domain: domain}, log.Default())
}

func t1TempDir() string {
	dir, _ := os.MkdirTemp("", "")
	return dir
}

func TestResolveEmail_EmptyDefaultsToDomain(t *testing.T) {
	service := newTestService("example.com")
	email, err := service.ResolveEmail("bob", "")
	assert.Nil(t, err)
	assert.Equal(t, "bob@example.com", email)
}

func TestResolveEmail_BlankDefaultsToDomain(t *testing.T) {
	service := newTestService("example.com")
	email, err := service.ResolveEmail("bob", "   ")
	assert.Nil(t, err)
	assert.Equal(t, "bob@example.com", email)
}

func TestResolveEmail_ValidKept(t *testing.T) {
	service := newTestService("example.com")
	email, err := service.ResolveEmail("bob", "bob@other.org")
	assert.Nil(t, err)
	assert.Equal(t, "bob@other.org", email)
}

func TestResolveEmail_InvalidRejected(t *testing.T) {
	service := newTestService("example.com")
	_, err := service.ResolveEmail("bob", "not-an-email")
	assert.NotNil(t, err)
}

func TestValidatePassword(t *testing.T) {
	assert.Nil(t, ValidatePassword("Password1"))
	assert.Nil(t, ValidatePassword("regularpass123"))
	assert.NotNil(t, ValidatePassword("short1"))
	assert.NotNil(t, ValidatePassword("alllettersonly"))
	assert.NotNil(t, ValidatePassword("12345678"))
}

func TestAddUser_WeakPasswordRejected(t *testing.T) {
	service := newTestService("example.com")
	err := service.AddUser("bob", "weak", "")
	assert.NotNil(t, err)
}

func TestAddUser_EmptyUsernameRejected(t *testing.T) {
	service := newTestService("example.com")
	err := service.AddUser("   ", "password", "")
	assert.NotNil(t, err)
}

func TestUserAttributes_MatchUsersApp(t *testing.T) {
	attrs := userAttributes("bob", "bob@example.com", 2001)
	assert.Equal(t, []string{"bob"}, attrs["cn"])
	assert.Equal(t, []string{"bob"}, attrs["sn"])
	assert.Equal(t, []string{"bob"}, attrs["givenName"])
	assert.Equal(t, []string{"bob"}, attrs["displayName"])
	assert.Equal(t, []string{"bob"}, attrs["uid"])
	assert.Equal(t, []string{"bob@example.com"}, attrs["mail"])
	assert.Equal(t, []string{"/home/bob"}, attrs["homeDirectory"])
	assert.Equal(t, []string{"/bin/bash"}, attrs["loginShell"])
	assert.Equal(t, []string{"2001"}, attrs["uidNumber"])
	assert.Equal(t, []string{"2001"}, attrs["gidNumber"])
	assert.Contains(t, attrs["objectClass"], "posixAccount")
	assert.Contains(t, attrs["objectClass"], "inetOrgPerson")
	assert.Contains(t, attrs["objectClass"], "person")
	assert.Contains(t, attrs["objectClass"], "simpleSecurityObject")
}

func TestInit(t *testing.T) {
	executor := &ExecutorStub{}
	ldap := New(&SnapServiceStub{}, t.TempDir(), t.TempDir(), t.TempDir(), executor, &PasswordChangerStub{}, DomainProviderStub{domain: "example.com"}, log.Default())
	err := ldap.Init()
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 1)
	assert.Contains(t, executor.executions[0], "slapadd.sh")
}

func TestApplyConfig_NotInstalled(t *testing.T) {
	executor := &ExecutorStub{}
	missing := path.Join(t.TempDir(), "missing")
	ldap := New(&SnapServiceStub{}, missing, t.TempDir(), t.TempDir(), executor, &PasswordChangerStub{}, DomainProviderStub{domain: "example.com"}, log.Default())
	err := ldap.ApplyConfig()
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 0)
}

func TestReset(t *testing.T) {
	executor := &ExecutorStub{}
	configDir := t.TempDir()
	err := os.MkdirAll(path.Join(configDir, "ldap"), os.ModePerm)
	assert.Nil(t, err)
	err = os.WriteFile(path.Join(configDir, "ldap", "init.ldif"), []byte("template"), 0644)
	assert.Nil(t, err)

	passwordChanger := &PasswordChangerStub{}
	ldap := New(&SnapServiceStub{}, t.TempDir(), t.TempDir(), configDir, executor, passwordChanger, DomainProviderStub{domain: "example.com"}, log.Default())
	err = ldap.Reset("name", "user", "password", "email")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 2)
	assert.Contains(t, executor.executions[0], "slapadd.sh")
	assert.Contains(t, executor.executions[1], "ldapadd.sh")
	assert.True(t, passwordChanger.changed)

}
