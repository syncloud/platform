package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syncloud/platform/storage/btrfs"
)

func btrfsCmd(userConfig *string, systemConfig *string) *cobra.Command {
	var cmdBtrfs = &cobra.Command{
		Use:   "btrfs",
		Short: "Show btrfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Init(*userConfig, *systemConfig)
			if err != nil {
				return err
			}
			return c.Call(func(btrfs *btrfs.Stats) error {
				info, err := btrfs.Info()
				if err != nil {
					return err
				}
				s, err := json.MarshalIndent(info, "", "\t")
				if err != nil {
					return err
				}
				fmt.Printf("btrfs info\n")
				fmt.Printf("%s\n", s)
				return nil
			})
		},
	}
	return cmdBtrfs
}
