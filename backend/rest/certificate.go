package rest

import (
	"fmt"
	"github.com/syncloud/platform/cert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/systemd"
	"net/http"
	"strings"
)

type Certificate struct {
	infoReader CertificateInfoReader
	journalCtl systemd.JournalCtlReader
}

type CertificateInfoReader interface {
	ReadCertificateInfo() *cert.Info
}

func NewCertificate(infoReader CertificateInfoReader, journalCtl systemd.JournalCtlReader) *Certificate {
	return &Certificate{
		infoReader: infoReader,
		journalCtl: journalCtl,
	}
}

func (c *Certificate) Certificate(_ *http.Request) (interface{}, error) {
	return c.infoReader.ReadCertificateInfo(), nil
}

func (c *Certificate) CertificateLog(_ *http.Request) (interface{}, error) {
	return c.journalCtl.ReadBackend(func(line string) bool {
		return strings.Contains(line, fmt.Sprintf(`"%s": "%s"`, log.CategoryKey, log.CategoryCertificate))
	}), nil
}
