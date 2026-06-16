package auth

import "strconv"

type UserAttributes struct{}

func NewUserAttributes() *UserAttributes {
	return &UserAttributes{}
}

func (a *UserAttributes) Build(username string, email string, id int) map[string][]string {
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
