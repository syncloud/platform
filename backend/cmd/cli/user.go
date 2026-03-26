package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/auth"
)

func userCmd(userConfig *string, systemConfig *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}
	cmd.AddCommand(
		userAddCmd(userConfig, systemConfig),
		userRemoveCmd(userConfig, systemConfig),
	)
	return cmd
}

func userAddCmd(userConfig *string, systemConfig *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <username>",
		Short: "Add a new user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			password, _ := cmd.Flags().GetString("password")
			if password == "" {
				return fmt.Errorf("--password is required")
			}
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(ldapService *auth.Service) {
				err := ldapService.AddUser(args[0], password)
				if err != nil {
					fmt.Printf("error: %v\n", err)
				} else {
					fmt.Printf("User '%s' added\n", args[0])
				}
			})
		},
	}
	cmd.Flags().String("password", "", "User password")
	return cmd
}

func userRemoveCmd(userConfig *string, systemConfig *string) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <username>",
		Short: "Remove a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(ldapService *auth.Service) {
				err := ldapService.RemoveUser(args[0])
				if err != nil {
					fmt.Printf("error: %v\n", err)
				} else {
					fmt.Printf("User '%s' removed\n", args[0])
				}
			})
		},
	}
}
