package auth

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestToLdapDc(t *testing.T) {
	assert.Equal(t, ToLdapDc("user.syncloud.it"), "dc=user,dc=syncloud,dc=it")
}

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

func (e *ExecutorStub) CommandOutput(name string, arg ...string) ([]byte, error) {
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

func TestMakeSecret(t *testing.T) {
	secret := makeSecret("syncloud")
	assert.Greater(t, len(secret), 1)
}

func TestInit(t *testing.T) {
	executor := &ExecutorStub{}
	ldap := New(&SnapServiceStub{}, createTempDir(), createTempDir(), createTempDir(), executor, &PasswordChangerStub{})
	err := ldap.Init()
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 1)
	assert.Contains(t, executor.executions[0], "slapadd.sh")
}

func TestReset(t *testing.T) {
	executor := &ExecutorStub{}
	configDir := createTempDir()
	err := os.MkdirAll(path.Join(configDir, "ldap"), os.ModePerm)
	assert.Nil(t, err)
	err = os.WriteFile(path.Join(configDir, "ldap", "init.ldif"), []byte("template"), 0644)
	assert.Nil(t, err)

	passwordChanger := &PasswordChangerStub{}
	ldap := New(&SnapServiceStub{}, createTempDir(), createTempDir(), configDir, executor, passwordChanger)
	err = ldap.Reset("name", "user", "password", "email")
	assert.Nil(t, err)
	assert.Len(t, executor.executions, 2)
	assert.Contains(t, executor.executions[0], "slapadd.sh")
	assert.Contains(t, executor.executions[1], "ldapadd.sh")
	assert.True(t, passwordChanger.changed)

}

func createTempDir() string {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		panic(err)
	}
	return dir
}
