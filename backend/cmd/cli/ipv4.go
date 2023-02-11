package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/ioc"
	"github.com/syncloud/platform/network"
	"os"
)

func ipv4Cmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdIpv4 = &cobra.Command{
		Use:   "ipv4 [public]",
		Short: "Print IPv4",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(iface *network.TcpInterfaces) {
				ip, err := iface.LocalIPv4()
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				fmt.Print(ip.String())
			})
		},
	}

	cmdIpv4.AddCommand(&cobra.Command{
		Use:   "public",
		Short: "Print public IPv4",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(iface *network.TcpInterfaces) {
				ip, err := iface.PublicIPv4()
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				fmt.Print(*ip)
			})
		},
	})
	return cmdIpv4
}
