package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/cron"
)

func cronCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdCron = &cobra.Command{
		Use:   "cron",
		Short: "Run cron job",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(cronService *cron.Cron) { cronService.StartSingle() })
		},
	}
	return cmdCron
}
