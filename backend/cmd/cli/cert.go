package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/cert"
)

func certCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdCert = &cobra.Command{
		Use:   "cert",
		Short: "Generate certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(certGenerator *cert.CertificateGenerator) error {
				return certGenerator.Generate()
			})
		},
	}
	return cmdCert
}
