package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/config"
)

func proxyCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdProxy = &cobra.Command{
		Use:   "proxy",
		Short: "Manage custom proxies",
	}

	var cmdProxyAdd = &cobra.Command{
		Use:   "add",
		Short: "Add a custom proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			host, _ := cmd.Flags().GetString("host")
			port, _ := cmd.Flags().GetInt("port")
			https, _ := cmd.Flags().GetBool("https")
			authelia, _ := cmd.Flags().GetBool("authelia")
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(configuration *config.UserConfig) {
				err := configuration.AddCustomProxy(name, host, port, https, authelia)
				if err != nil {
					fmt.Printf("error: %s\n", err)
				}
			})
		},
	}
	cmdProxyAdd.Flags().String("name", "", "proxy name")
	cmdProxyAdd.Flags().String("host", "", "backend host")
	cmdProxyAdd.Flags().Int("port", 0, "backend port")
	cmdProxyAdd.Flags().Bool("https", false, "use https for backend")
	cmdProxyAdd.Flags().Bool("authelia", false, "protect with Authelia (SSO + 2FA)")
	_ = cmdProxyAdd.MarkFlagRequired("name")
	_ = cmdProxyAdd.MarkFlagRequired("host")
	_ = cmdProxyAdd.MarkFlagRequired("port")
	cmdProxy.AddCommand(cmdProxyAdd)

	var cmdProxyRemove = &cobra.Command{
		Use:   "remove",
		Short: "Remove a custom proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(configuration *config.UserConfig) {
				err := configuration.RemoveCustomProxy(name)
				if err != nil {
					fmt.Printf("error: %s\n", err)
				}
			})
		},
	}
	cmdProxyRemove.Flags().String("name", "", "proxy name")
	_ = cmdProxyRemove.MarkFlagRequired("name")
	cmdProxy.AddCommand(cmdProxyRemove)

	var cmdProxyList = &cobra.Command{
		Use:   "list",
		Short: "List custom proxies",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(configuration *config.UserConfig) {
				proxies, err := configuration.CustomProxies()
				if err != nil {
					fmt.Printf("error: %s\n", err)
					return
				}
				data, _ := json.Marshal(proxies)
				fmt.Println(string(data))
			})
		},
	}
	cmdProxy.AddCommand(cmdProxyList)

	return cmdProxy
}
