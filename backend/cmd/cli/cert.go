package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/cert"
)

func certCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var fake bool
	var cmdCert = &cobra.Command{
		Use:   "cert",
		Short: "Generate certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			if fake {
				return c.Call(func(fakeGenerator *cert.Fake) error {
					return fakeGenerator.Generate()
				})
			}
			return c.Call(func(certGenerator *cert.CertificateGenerator) error {
				return certGenerator.Generate()
			})
		},
	}
	cmdCert.Flags().BoolVar(&fake, "fake", false, "force fake (self-signed) CA + cert pair regeneration, bypassing the Subject==domain skip")
	return cmdCert
}
