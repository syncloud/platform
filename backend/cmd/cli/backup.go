package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/backup"
)

func backupCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup create/restore",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create [app]",
		Short: "Backup",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(backup *backup.Backup) error {
				app := args[0]
				return backup.Create(app)
			})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "restore [file]",
		Short: "Restore file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(backup *backup.Backup) error {
				return backup.Restore(args[0])
			})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(backup *backup.Backup) error {
				list, err := backup.List()
				if err != nil {
					return err
				}
				s, err := json.MarshalIndent(list, "", "\t")
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", s)
				return nil
			})
		},
	})
	return cmd
}
