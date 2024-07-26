package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhangx1n/xim/client"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: clientHandler,
}

func clientHandler(cmd *cobra.Command, args []string) {
	client.RunMain()
}
