package cert

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
)

func readCertificate(t *testing.T, path string) *x509.Certificate {
	raw, err := os.ReadFile(path)
	assert.NoError(t, err)
	block, _ := pem.Decode(raw)
	assert.NotNil(t, block)
	certificate, err := x509.ParseCertificate(block.Bytes)
	assert.NoError(t, err)
	return certificate
}

func generateFake(t *testing.T) (*x509.Certificate, *x509.Certificate) {
	dir := t.TempDir()
	config := &GeneratorSystemConfigStub{
		certFile:   dir + "/cert.pem",
		keyFile:    dir + "/key.pem",
		caSertFile: dir + "/ca.pem",
		caKeyFile:  dir + "/ca.key.pem",
	}
	fake := NewFake(
		config,
		&GeneratorUserConfigStub{domain: "example.com"},
		&ProviderStub{now: time.Now()},
		SubjectOrganization,
		DefaultDuration,
		log.Default(),
	)
	assert.NoError(t, fake.Generate())
	return readCertificate(t, config.caSertFile), readCertificate(t, config.certFile)
}

func TestFakeCertificateHasAuthorityKeyId(t *testing.T) {
	ca, leaf := generateFake(t)

	assert.NotEmpty(t, ca.SubjectKeyId, "the leaf takes its authority key id from this")
	assert.NotEmpty(t, leaf.AuthorityKeyId,
		"python 3.13 onwards verifies with VERIFY_X509_STRICT by default, which rejects "+
			"a certificate carrying no authority key identifier")
	assert.Equal(t, ca.SubjectKeyId, leaf.AuthorityKeyId)
}

func TestFakeCertificateChainVerifies(t *testing.T) {
	ca, leaf := generateFake(t)

	roots := x509.NewCertPool()
	roots.AddCert(ca)

	_, err := leaf.Verify(x509.VerifyOptions{
		Roots:     roots,
		DNSName:   "app.example.com",
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
	assert.NoError(t, err)
}
