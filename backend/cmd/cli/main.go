package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/config"
)

func main() {
	var rootCmd = &cobra.Command{Use: "cli"}
	userConfig := rootCmd.PersistentFlags().String("user-config", config.DefaultConfigDb, "user config sqlite db")
	systemConfig := rootCmd.PersistentFlags().String("system-config", config.DefaultSystemConfig, "system config file")

	rootCmd.AddCommand(
		ipv4Cmd(userConfig, systemConfig),
		ipv6Cmd(userConfig, systemConfig),
		configCmd(userConfig, systemConfig),
		cronCmd(userConfig, systemConfig),
		certCmd(userConfig, systemConfig),
		btrfsCmd(userConfig, systemConfig),
		backupCmd(userConfig, systemConfig),
	)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
