package auth

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegexp = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type DomainProvider interface {
	GetDeviceDomain() string
}

type EmailResolver struct {
	domain DomainProvider
}

func NewEmailResolver(domain DomainProvider) *EmailResolver {
	return &EmailResolver{domain: domain}
}

func (r *EmailResolver) Resolve(username string, email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Sprintf("%s@%s", username, r.domain.GetDeviceDomain()), nil
	}
	if !emailRegexp.MatchString(email) {
		return "", fmt.Errorf("invalid email address: %s", email)
	}
	return email, nil
}
