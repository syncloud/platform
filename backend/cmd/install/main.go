package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/hook"
	"github.com/syncloud/platform/ioc"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "install",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := ioc.Init(config.DefaultConfigDb, config.DefaultSystemConfig, backup.Dir, backup.VarDir)
			if err != nil {
				return err
			}
			return c.Call(func(install *hook.Install) error {
				return install.Run()
			})
		},
	}

	err := rootCmd.Execute()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
