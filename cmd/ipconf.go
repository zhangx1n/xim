package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhangx1n/xim/ipconf"
)

func init() {
	rootCmd.AddCommand(ipConfCmd)
}

var ipConfCmd = &cobra.Command{
	Use: "ipconf",
	Run: IpConfHandler,
}

func IpConfHandler(cmd *cobra.Command, args []string) {
	ipconf.RunMain(ConfigPath)
}
