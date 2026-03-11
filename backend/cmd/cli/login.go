package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/syncloud/platform/auth"
)

const TokenFile = "/var/snap/platform/common/login-token"

func loginCmd(userConfig *string, systemConfig *string) *cobra.Command {
	return &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Generate a one-time login token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			password := args[1]

			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			var authErr error
			err = c.Call(func(authService *auth.Service) {
				ok, err := authService.Authenticate(username, password)
				if err != nil {
					authErr = fmt.Errorf("authentication failed: %w", err)
					return
				}
				if !ok {
					authErr = fmt.Errorf("authentication failed")
					return
				}

				tokenBytes := make([]byte, 32)
				_, err = rand.Read(tokenBytes)
				if err != nil {
					authErr = fmt.Errorf("unable to generate token: %w", err)
					return
				}
				token := hex.EncodeToString(tokenBytes)

				err = os.WriteFile(TokenFile, []byte(token+":"+username), 0600)
				if err != nil {
					authErr = fmt.Errorf("unable to write token file: %w", err)
					return
				}

				fmt.Print(token)
			})
			if err != nil {
				return err
			}
			return authErr
		},
	}
}
