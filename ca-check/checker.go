package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

type Checker struct {
	roots    *x509.CertPool
	deadline time.Time
}

func NewChecker(dir string, deadline time.Time) (*Checker, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.crt"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no certificates in %s", dir)
	}
	roots := x509.NewCertPool()
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if !roots.AppendCertsFromPEM(content) {
			return nil, fmt.Errorf("invalid certificate %s", file)
		}
	}
	return &Checker{roots: roots, deadline: deadline}, nil
}

func (c *Checker) Check(host string) error {
	chains, err := c.chains(host)
	if err != nil {
		return err
	}
	for _, chain := range chains {
		expiry := chain[len(chain)-1].NotAfter
		for _, certificate := range chain[min(2, len(chain)-1):] {
			if certificate.NotAfter.Before(expiry) {
				expiry = certificate.NotAfter
			}
		}
		fmt.Printf("%s: %s, expires %s\n", host, describe(chain), expiry.Format(time.DateOnly))
		if expiry.After(c.deadline) {
			return nil
		}
	}
	return fmt.Errorf("%s: every chain to a shipped root expires before %s", host, c.deadline.Format(time.DateOnly))
}

func (c *Checker) chains(host string) ([][]*x509.Certificate, error) {
	config := &tls.Config{RootCAs: c.roots, ServerName: host}
	dialer := &net.Dialer{Timeout: 30 * time.Second}
	var err error
	for attempt := 1; attempt <= 5; attempt++ {
		var conn *tls.Conn
		conn, err = tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, "443"), config)
		if err == nil {
			defer conn.Close()
			return conn.ConnectionState().VerifiedChains, nil
		}
		var verification *tls.CertificateVerificationError
		if errors.As(err, &verification) {
			return nil, err
		}
		fmt.Printf("%s: attempt %d: %v\n", host, attempt, err)
		time.Sleep(time.Duration(attempt) * 5 * time.Second)
	}
	return nil, err
}

func describe(chain []*x509.Certificate) string {
	description := ""
	for i, certificate := range chain {
		if i > 0 {
			description += " -> "
		}
		description += certificate.Subject.CommonName
	}
	return description
}
