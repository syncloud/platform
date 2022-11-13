package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/ioc"
	"github.com/syncloud/platform/storage/btrfs"
	"os"
)

func btrfsCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdBtrfs = &cobra.Command{
		Use:   "btrfs",
		Short: "Show btrfs",
		Run: func(cmd *cobra.Command, args []string) {
			Init(*userConfig, *systemConfig)
			ioc.Call(func(btrfs *btrfs.Stats) {
				info, err := btrfs.Info()
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				s, err := json.MarshalIndent(info, "", "\t")
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}
				fmt.Printf("btrfs info\n")
				fmt.Printf("%s\n", s)
			})
		},
	}
	return cmdBtrfs
}
