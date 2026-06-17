package auth

import (
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
)

type UserBuilder struct {
	passwordHasher *PasswordHasher
}

func NewUserBuilder(passwordHasher *PasswordHasher) *UserBuilder {
	return &UserBuilder{passwordHasher: passwordHasher}
}

func (b *UserBuilder) Build(username string, email string, id int, password string) *ldap.AddRequest {
	idStr := strconv.Itoa(id)
	req := ldap.NewAddRequest(fmt.Sprintf("cn=%s,%s", username, UsersDn), nil)
	req.Attribute("objectClass", []string{"person", "inetOrgPerson", "posixAccount", "simpleSecurityObject"})
	req.Attribute("cn", []string{username})
	req.Attribute("sn", []string{username})
	req.Attribute("givenName", []string{username})
	req.Attribute("displayName", []string{username})
	req.Attribute("uid", []string{username})
	req.Attribute("uidNumber", []string{idStr})
	req.Attribute("gidNumber", []string{idStr})
	req.Attribute("homeDirectory", []string{"/home/" + username})
	req.Attribute("loginShell", []string{"/bin/bash"})
	req.Attribute("mail", []string{email})
	req.Attribute("userPassword", []string{b.passwordHasher.Hash(password)})
	return req
}
