package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/redirect"
	"log"
	"os"

	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/storage"
)

func main() {

	var rootCmd = &cobra.Command{Use: "backend"}
	defaultConfigDb := fmt.Sprintf("%s/platform.db", os.Getenv("SNAP_COMMON"))
	configDb := rootCmd.PersistentFlags().String("config", defaultConfigDb, "sqlite config db")
	redirectDomain := rootCmd.PersistentFlags().String("redirect-domain", "syncloud.it", "redirect domain")
	redirectUrl := rootCmd.PersistentFlags().String("redirect-url", "https://api.syncloud.it", "redirect url")

	var tcpCmd = &cobra.Command{
		Use:   "tcp [address]",
		Short: "listen on a tcp address, like localhost:8080",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Backend(*configDb, *redirectDomain, *redirectUrl).Start("tcp", args[0])
		},
	}

	var unixSocketCmd = &cobra.Command{
		Use:   "unix [address]",
		Args:  cobra.ExactArgs(1),
		Short: "listen on a unix socket, like /tmp/backend.sock",
		Run: func(cmd *cobra.Command, args []string) {
			_ = os.Remove(args[0])
			Backend(*configDb, *redirectDomain, *redirectUrl).Start("unix", args[0])
		},
	}

	rootCmd.AddCommand(tcpCmd, unixSocketCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Print("error: ", err)
		os.Exit(1)
	}
}

func Backend(configDb string, redirectDomain string, redirectUrl string) *rest.Backend {
	master := job.NewMaster()
	backupService := backup.NewDefault()
	eventTrigger := event.New()
	installerService := installer.New()
	storageService := storage.New()
	oldConfig := fmt.Sprintf("%s/user_platform.cfg", os.Getenv("SNAP_COMMON"))
	configuration := config.New(configDb, oldConfig, redirectDomain, redirectUrl)
	redirectService := redirect.New(configuration)
	worker := job.NewWorker(master)
	return rest.NewBackend(master, backupService, eventTrigger, worker, redirectService, installerService, storageService)

}
