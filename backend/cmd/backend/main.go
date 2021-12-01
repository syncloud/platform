package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/activation"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/certificate/certbot"
	"github.com/syncloud/platform/certificate/fake"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/connection"
	"github.com/syncloud/platform/cron"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/logger"
	"github.com/syncloud/platform/nginx"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/snap"
	"github.com/syncloud/platform/systemd"
	"log"
	"os"
	"time"

	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/installer"
	"github.com/syncloud/platform/job"
	"github.com/syncloud/platform/rest"
	"github.com/syncloud/platform/storage"
)

func main() {

	log.SetFlags(0)
	log.SetOutput(&logger.Logger{})

	var rootCmd = &cobra.Command{Use: "backend"}
	configDb := rootCmd.PersistentFlags().String("config", config.DefaultConfigDb, "sqlite config db")

	var tcpCmd = &cobra.Command{
		Use:   "tcp [address]",
		Short: "listen on a tcp address, like localhost:8080",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			backend, err := Backend(*configDb)
			if err != nil {
				log.Print("error: ", err)
				os.Exit(1)
			}
			backend.Start("tcp", args[0])
		},
	}

	var unixSocketCmd = &cobra.Command{
		Use:   "unix [address]",
		Args:  cobra.ExactArgs(1),
		Short: "listen on a unix socket, like /tmp/backend.sock",
		Run: func(cmd *cobra.Command, args []string) {
			_ = os.Remove(args[0])
			backend, err := Backend(*configDb)
			if err != nil {
				log.Print("error: ", err)
				os.Exit(1)
			}
			backend.Start("unix", args[0])
		},
	}

	rootCmd.AddCommand(tcpCmd, unixSocketCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Print("error: ", err)
		os.Exit(1)
	}
}

func Backend(configDb string) (*rest.Backend, error) {

	cronService := cron.New(cron.Job, time.Minute*5)
	cronService.Start()
	master := job.NewMaster()
	backupService := backup.NewDefault()
	snapClient := snap.NewClient()
	snapd := snap.New(snapClient)
	eventTrigger := event.New(snapd)
	installerService := installer.New()
	storageService := storage.New()
	userConfig, err := config.NewUserConfig(configDb, config.OldConfig)
	if err != nil {
		return nil, err
	}
	id := identification.New()
	redirectService := redirect.New(userConfig, id)
	worker := job.NewWorker(master)
	systemConfig, err := config.NewSystemConfig(config.File)
	if err != nil {
		return nil, err
	}
	snapService := snap.NewService()
	dataDir := systemConfig.DataDir()
	appDir := systemConfig.AppDir()
	configDir := systemConfig.ConfigDir()
	ldapService := auth.New(snapService, dataDir, appDir, configDir)
	nginxService := nginx.New(systemd.New(), systemConfig, userConfig)
	device := activation.NewDevice(userConfig, ldapService, nginxService, eventTrigger)
	internetChecker := connection.NewInternetChecker()
	realCert := certbot.New(redirectService, userConfig, systemConfig)
	fakeCert := fake.New(systemConfig)
	activationManaged := activation.NewManaged(internetChecker, userConfig, redirectService, device, realCert, fakeCert)
	activationCustom := activation.NewCustom(internetChecker, userConfig, redirectService, device, fakeCert)
	activate := rest.NewActivateBackend(activationManaged, activationCustom)
	backend := rest.NewBackend(master, backupService, eventTrigger, worker, redirectService,
		installerService, storageService, id, activate, userConfig)
	return backend, nil

}
