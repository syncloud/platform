package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/cert"
	"github.com/syncloud/platform/ioc"
	"os"
)

func certCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdCert = &cobra.Command{
		Use:   "cert",
		Short: "Generate certificate",
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(certGenerator *cert.CertificateGenerator) {
				err := certGenerator.Generate()
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
			})
		},
	}
	return cmdCert
}
