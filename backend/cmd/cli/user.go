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
			email, _ := cmd.Flags().GetString("email")
			admin, _ := cmd.Flags().GetBool("admin")
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(userManager *auth.UserManager) {
				err := userManager.AddUser(args[0], password, email, admin)
				if err != nil {
					fmt.Printf("error: %v\n", err)
				} else {
					fmt.Printf("User '%s' added\n", args[0])
				}
			})
		},
	}
	cmd.Flags().String("password", "", "User password")
	cmd.Flags().String("email", "", "User email (defaults to <username>@<device domain>)")
	cmd.Flags().Bool("admin", false, "Make the user an admin")
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
			return c.Call(func(userManager *auth.UserManager) {
				err := userManager.RemoveUser(args[0])
				if err != nil {
					fmt.Printf("error: %v\n", err)
				} else {
					fmt.Printf("User '%s' removed\n", args[0])
				}
			})
		},
	}
}
