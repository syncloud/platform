package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/config"
)

func configCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var configFile string
	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Manage config",
	}
	cmdConfig.PersistentFlags().StringVar(&configFile, "file", config.DefaultConfigDb, "config file")

	var cmdConfigSet = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set config key value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(db *config.Db) {
				key := args[0]
				value := args[1]
				db.Upsert(key, value)
				fmt.Printf("set config: %s, key: %s, value: %s\n", configFile, key, value)
			})
		},
	}
	cmdConfig.AddCommand(cmdConfigSet)

	var cmdConfigGet = &cobra.Command{
		Use:   "get [key]",
		Short: "Get config key value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(db *config.Db) {
				fmt.Println(db.Get(args[0], ""))
			})
		},
	}
	cmdConfig.AddCommand(cmdConfigGet)

	var cmdConfigList = &cobra.Command{
		Use:   "list",
		Short: "List config key value",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(db *config.Db) {
				for key, value := range db.List() {
					fmt.Printf("%s:%s\n", key, value)
				}
			})
		},
	}
	cmdConfig.AddCommand(cmdConfigList)
	return cmdConfig
}
