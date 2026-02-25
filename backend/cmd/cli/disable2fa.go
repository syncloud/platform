package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/config"
)

func disable2faCmd(userConfig *string, systemConfig *string) *cobra.Command {
	return &cobra.Command{
		Use:   "disable-2fa",
		Short: "Disable two-factor authentication",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(configuration *config.UserConfig, authelia *auth.Authelia) {
				configuration.SetTwoFactorEnabled(false)
				fmt.Println("2FA disabled in config")
				err := authelia.InitConfig()
				if err != nil {
					fmt.Printf("warning: unable to regenerate authelia config: %v\n", err)
				} else {
					fmt.Println("Authelia config regenerated and service restarted")
				}
			})
		},
	}
}
