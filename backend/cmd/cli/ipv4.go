package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/network"
)

func ipv4Cmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdIpv4 = &cobra.Command{
		Use:   "ipv4 [public]",
		Short: "Print IPv4",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				panic(err)
			}
			return c.Call(func(iface *network.TcpInterfaces) error {
				ip, err := iface.LocalIPv4()
				if err != nil {
					return err
				}
				fmt.Print(ip.String())
				return nil
			})
		},
	}

	cmdIpv4.AddCommand(&cobra.Command{
		Use:   "public",
		Short: "Print public IPv4",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				panic(err)
			}
			return c.Call(func(iface *network.TcpInterfaces) error {
				ip, err := iface.PublicIPv4()
				if err != nil {
					return err
				}
				fmt.Print(*ip)
				return nil
			})
		},
	})
	return cmdIpv4
}
