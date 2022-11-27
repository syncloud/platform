package rest

import (
	"fmt"
	"github.com/syncloud/platform/cert"
	"github.com/syncloud/platform/log"
	"net/http"
	"strings"
)

type Journal interface {
	ReadBackend(predicate func(string) bool) []string
}

type Certificate struct {
	infoReader CertificateInfoReader
	journal    Journal
}

type CertificateInfoReader interface {
	ReadCertificateInfo() *cert.Info
}

func NewCertificate(infoReader CertificateInfoReader, journal Journal) *Certificate {
	return &Certificate{
		infoReader: infoReader,
		journal:    journal,
	}
}

func (c *Certificate) Certificate(_ *http.Request) (interface{}, error) {
	return c.infoReader.ReadCertificateInfo(), nil
}

func (c *Certificate) CertificateLog(_ *http.Request) (interface{}, error) {
	return c.journal.ReadBackend(func(line string) bool {
		return strings.Contains(line, fmt.Sprintf(`"%s": "%s"`, log.CategoryKey, log.CategoryCertificate))
	}), nil
}
