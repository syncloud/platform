package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/cron"
	"github.com/syncloud/platform/ioc"
	"github.com/syncloud/platform/rest"
	"log"
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
		Run: func(cmd *cobra.Command, args []string) {
			Start(*userConfig, *systemConfig, "tcp", args[0])
		},
	}

	var unixSocketCmd = &cobra.Command{
		Use:   "unix [address]",
		Args:  cobra.ExactArgs(1),
		Short: "listen on a unix socket, like /tmp/backend.sock",
		Run: func(cmd *cobra.Command, args []string) {
			_ = os.Remove(args[0])
			Start(*userConfig, *systemConfig, "unix", args[0])
		},
	}

	rootCmd.AddCommand(tcpCmd, unixSocketCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Print("error: ", err)
		os.Exit(1)
	}
}

func Start(userConfig string, systemConfig string, socketType string, socket string) {
	ioc.Init(userConfig, systemConfig, backup.Dir, backup.VarDir)
	ioc.Call(func(cronService *cron.Cron) { cronService.StartScheduler() })
	ioc.Call(func(backupService *backup.Backup) { backupService.Init() })
	ioc.Call(func(backend *rest.Backend) error { return backend.Start(socketType, socket) })
}
