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

func newTestService(domain string) *Service {
	return New(&SnapServiceStub{}, t1TempDir(), t1TempDir(), t1TempDir(), &ExecutorStub{}, NewLdapClient(), &PasswordChangerStub{}, NewPasswordValidator(), NewPasswordHasher(), NewEmailResolver(DomainProviderStub{domain: domain}), NewUserBuilder(NewPasswordHasher()), log.Default())
}

func t1TempDir() string {
	dir, _ := os.MkdirTemp("", "")
	return dir
}

func TestAddUser_WeakPasswordRejected(t *testing.T) {
	service := newTestService("example.com")
	err := service.AddUser("bob", "weak", "", false)
	assert.NotNil(t, err)
}

func TestAddUser_EmptyUsernameRejected(t *testing.T) {
	service := newTestService("example.com")
	err := service.AddUser("   ", "password", "", false)
	assert.NotNil(t, err)
}

func TestInit(t *testing.T) {
	executor := &ExecutorStub{}
	ldap := New(&SnapServiceStub{}, t.TempDir(), t.TempDir(), t.TempDir(), executor, NewLdapClient(), &PasswordChangerStub{}, NewPasswordValidator(), NewPasswordHasher(), NewEmailResolver(DomainProviderStub{domain: "example.com"}), NewUserBuilder(NewPasswordHasher()), log.Default())
	err := ldap.Init()
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 1)
	assert.Contains(t, executor.executions[0], "slapadd.sh")
}

func TestApplyConfig_NotInstalled(t *testing.T) {
	executor := &ExecutorStub{}
	missing := path.Join(t.TempDir(), "missing")
	ldap := New(&SnapServiceStub{}, missing, t.TempDir(), t.TempDir(), executor, NewLdapClient(), &PasswordChangerStub{}, NewPasswordValidator(), NewPasswordHasher(), NewEmailResolver(DomainProviderStub{domain: "example.com"}), NewUserBuilder(NewPasswordHasher()), log.Default())
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
	ldap := New(&SnapServiceStub{}, t.TempDir(), t.TempDir(), configDir, executor, NewLdapClient(), passwordChanger, NewPasswordValidator(), NewPasswordHasher(), NewEmailResolver(DomainProviderStub{domain: "example.com"}), NewUserBuilder(NewPasswordHasher()), log.Default())
	err = ldap.Reset("name", "user", "password", "email")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 2)
	assert.Contains(t, executor.executions[0], "slapadd.sh")
	assert.Contains(t, executor.executions[1], "ldapadd.sh")
	assert.True(t, passwordChanger.changed)

}
