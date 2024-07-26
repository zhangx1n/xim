package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhangx1n/xim/gateway"
)

func init() {
	rootCmd.AddCommand(gatewayCmd)

}

var gatewayCmd = &cobra.Command{
	Use: "gateway",
	Run: GatewayHandler,
}

func GatewayHandler(cmd *cobra.Command, args []string) {
	gateway.RunMain(ConfigPath)
}
