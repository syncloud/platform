package auth

import (
	"fmt"
	"unicode"
)

const passwordMinLength = 8

type PasswordValidator struct {
	minLength int
}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{minLength: passwordMinLength}
}

func (v *PasswordValidator) Validate(password string) error {
	if len(password) < v.minLength {
		return fmt.Errorf("password must be at least %d characters", v.minLength)
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
