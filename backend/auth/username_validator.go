package auth

import (
	"fmt"
	"regexp"
)

var usernameRegexp = regexp.MustCompile(`^[a-z][a-z0-9._-]{1,31}$`)

type UsernameValidator struct {
}

func NewUsernameValidator() *UsernameValidator {
	return &UsernameValidator{}
}

func (v *UsernameValidator) Validate(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if !usernameRegexp.MatchString(username) {
		return fmt.Errorf("username must be 2-32 characters, start with a lowercase letter and contain only lowercase letters, numbers, dots, dashes and underscores")
	}
	return nil
}
