package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/cron"
	"github.com/syncloud/platform/ioc"
)

func cronCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdCron = &cobra.Command{
		Use:   "cron",
		Short: "Run cron job",
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(cronService *cron.Cron) { cronService.StartSingle() })
		},
	}
	return cmdCron
}
