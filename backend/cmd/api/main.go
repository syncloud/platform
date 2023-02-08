package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/ioc"
	"os"
)

func main() {

	var rootCmd = &cobra.Command{Use: "api"}
	configDb := rootCmd.PersistentFlags().String("config", config.DefaultConfigDb, "sqlite config db")

	var unixSocketCmd = &cobra.Command{
		Use:   "unix [address]",
		Args:  cobra.ExactArgs(1),
		Short: "listen on a unix socket, like /tmp/api.sock",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = os.Remove(args[0])
			ioc.InitInternalApi(*configDb, config.DefaultSystemConfig, backup.Dir, backup.VarDir, "unix", args[0])
			return ioc.Start()
		},
	}

	rootCmd.AddCommand(unixSocketCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
