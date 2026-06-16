package auth

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type LdapClient struct {
	url string
}

func NewLdapClient() *LdapClient {
	return &LdapClient{url: "ldap://localhost:389"}
}

func (c *LdapClient) Connect() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(c.url)
	if err != nil {
		return nil, fmt.Errorf("ldap connect: %w", err)
	}
	if err := conn.Bind(Domain, "syncloud"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ldap root bind: %w", err)
	}
	return conn, nil
}

func (c *LdapClient) Disconnect(conn *ldap.Conn) {
	conn.Close()
}
