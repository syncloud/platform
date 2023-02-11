package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/ioc"
	"github.com/syncloud/platform/network"
	"net"
	"os"
)

func ipv6Cmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdIpv6 = &cobra.Command{
		Use:   "ipv6 [prefix]",
		Short: "Print IPv6",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(iface *network.TcpInterfaces) {
				ip, err := iface.IPv6Addr()
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				fmt.Print(ip.String())
			})
		},
	}
	var prefixSize int
	var cmdIpv6prefix = &cobra.Command{
		Use:   "prefix",
		Short: "Print IPv6 prefix",
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(iface *network.TcpInterfaces) {
				ip, err := iface.IPv6Addr()
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				fmt.Printf("%v/%v", ip.Mask(net.CIDRMask(prefixSize, 128)), prefixSize)
			})
		},
	}
	cmdIpv6prefix.Flags().IntVarP(&prefixSize, "size", "s", 64, "Prefix size")
	cmdIpv6.AddCommand(cmdIpv6prefix)
	return cmdIpv6
}
