package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/ioc"
	"os"
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
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(backup *backup.Backup) {
				app := args[0]
				err := backup.Create(app)
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
			})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "restore [file]",
		Short: "Restore file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(backup *backup.Backup) {
				err := backup.Restore(args[0])
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
			})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List backups",
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(backup *backup.Backup) {
				list, err := backup.List()
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				s, err := json.MarshalIndent(list, "", "\t")
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				fmt.Printf("%s\n", s)
			})
		},
	})
	return cmd
}
