package auth

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"
)

type Authenticator struct {
	ldapClient *LdapClient
	logger     *zap.Logger
}

func NewAuthenticator(ldapClient *LdapClient, logger *zap.Logger) *Authenticator {
	return &Authenticator{ldapClient: ldapClient, logger: logger}
}

func (a *Authenticator) Authenticate(username string, password string) (bool, error) {
	conn, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		return false, err
	}
	defer a.ldapClient.Disconnect(conn)
	err = conn.Bind(fmt.Sprintf("cn=%s,ou=users,dc=syncloud,dc=org", username), password)
	if err != nil {
		a.logger.Error("ldap error", zap.Error(err))
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
