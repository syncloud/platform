package cron

import (
	"github.com/syncloud/platform/cert"
)

type CertificateJob struct {
	certGenerator cert.Generator
}

func NewCertificateJob(certGenerator cert.Generator) *CertificateJob {
	return &CertificateJob{
		certGenerator: certGenerator,
	}
}

func (j *CertificateJob) Run() error {
	return j.certGenerator.Generate()
}
