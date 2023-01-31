package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/backup"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/ioc"
	"github.com/syncloud/platform/rest"
	"log"
	"os"
)

func main() {

	var rootCmd = &cobra.Command{Use: "api"}
	configDb := rootCmd.PersistentFlags().String("config", config.DefaultConfigDb, "sqlite config db")

	var unixSocketCmd = &cobra.Command{
		Use:   "unix [address]",
		Args:  cobra.ExactArgs(1),
		Short: "listen on a unix socket, like /tmp/api.sock",
		Run: func(cmd *cobra.Command, args []string) {
			_ = os.Remove(args[0])
			ioc.Init(*configDb, config.DefaultSystemConfig, backup.Dir, backup.VarDir)
			ioc.Call(func(api *rest.Api) { api.Start("unix", args[0]) })

		},
	}

	rootCmd.AddCommand(unixSocketCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Print("error: ", err)
		os.Exit(1)
	}
}
