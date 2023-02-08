package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/ioc"
	"os"
)

func main() {

	var rootCmd = &cobra.Command{Use: "backend"}
	userConfig := rootCmd.PersistentFlags().String("user-config", config.DefaultConfigDb, "sqlite config db")
	systemConfig := rootCmd.PersistentFlags().String("system-config", config.DefaultSystemConfig, "system config")

	var tcpCmd = &cobra.Command{
		Use:   "tcp [address]",
		Short: "listen on a tcp address, like localhost:8080",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ioc.InitPublicApi(*userConfig, *systemConfig, backup.Dir, backup.VarDir, "tcp", args[0])
			return ioc.Start()
		},
	}

	var unixSocketCmd = &cobra.Command{
		Use:   "unix [address]",
		Args:  cobra.ExactArgs(1),
		Short: "listen on a unix socket, like /tmp/backend.sock",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = os.Remove(args[0])
			ioc.InitPublicApi(*userConfig, *systemConfig, backup.Dir, backup.VarDir, "unix", args[0])
			return ioc.Start()
		},
	}

	rootCmd.AddCommand(tcpCmd, unixSocketCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
