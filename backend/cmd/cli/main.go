package main

import (
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/cron"
	"github.com/syncloud/platform/logger"
	"github.com/syncloud/platform/network"
	"log"
	"net"
)

func main() {

	log.SetFlags(0)
	log.SetOutput(&logger.Logger{})

	var cmdIpv4 = &cobra.Command{
		Use:   "ipv4 [public]",
		Short: "Print IPv4",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := network.LocalIPv4()
			if err != nil {
				log.Fatal(err)
			}
			log.Print(ip.String())
		},
	}
	var cmdIpv4public = &cobra.Command{
		Use:   "public",
		Short: "Print public IPv4",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := network.PublicIPv4()
			if err != nil {
				log.Fatal(err)
			}
			log.Print(ip)
		},
	}
	cmdIpv4.AddCommand(cmdIpv4public)

	var cmdIpv6 = &cobra.Command{
		Use:   "ipv6 [prefix]",
		Short: "Print IPv6",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := network.IPv6()
			if err != nil {
				log.Fatal(err)
			}
			log.Print(ip.String())
		},
	}
	var prefixSize int
	var cmdIpv6prefix = &cobra.Command{
		Use:   "prefix",
		Short: "Print IPv6 prefix",
		//Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := network.IPv6()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%v/%v", ip.Mask(net.CIDRMask(prefixSize, 128)), prefixSize)
		},
	}
	cmdIpv6prefix.Flags().IntVarP(&prefixSize, "size", "s", 64, "Prefix size")
	cmdIpv6.AddCommand(cmdIpv6prefix)

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
		Run: func(cmd *cobra.Command, args []string) {
			configuration, err := config.NewUserConfig(configFile, config.OldConfig, "")
			if err != nil {
				log.Fatal(err)
			}
			key := args[0]
			value := args[1]
			configuration.Upsert(key, value)
			log.Printf("set config: %s, key: %s, value: %s\n", configFile, key, value)
		},
	}
	cmdConfig.AddCommand(cmdConfigSet)

	var cmdConfigGet = &cobra.Command{
		Use:   "get [key]",
		Short: "Get config key value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configuration, err := config.NewUserConfig(configFile, config.OldConfig, "")
			if err != nil {
				log.Fatal(err)
			}
			log.Println(configuration.Get(args[0], ""))
		},
	}
	cmdConfig.AddCommand(cmdConfigGet)

	var cmdConfigList = &cobra.Command{
		Use:   "list",
		Short: "List config key value",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			configuration, err := config.NewUserConfig(configFile, config.OldConfig, "")
			if err != nil {
				log.Fatal(err)
			}
			for key, value := range configuration.List() {
				log.Printf("%s:%s\n", key, value)
			}
		},
	}
	cmdConfig.AddCommand(cmdConfigList)

	var cmdCron = &cobra.Command{
		Use:   "cron",
		Short: "Run cron job",
		Run: func(cmd *cobra.Command, args []string) {
			err := cron.Job()
			if err != nil {
				log.Fatalf("error: %s\n", err)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "cli"}
	rootCmd.AddCommand(cmdIpv4, cmdIpv6, cmdConfig, cmdCron)
	rootCmd.Execute()
}
