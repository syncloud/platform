package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/auth"
)

func resetTotpCmd(userConfig *string, systemConfig *string) *cobra.Command {
	return &cobra.Command{
		Use:   "reset-totp <username>",
		Short: "Delete TOTP registration for a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(authelia *auth.Authelia) {
				err := authelia.ResetTOTP(args[0])
				if err != nil {
					fmt.Printf("error: %v\n", err)
				} else {
					fmt.Println("TOTP registration deleted")
				}
			})
		},
	}
}
