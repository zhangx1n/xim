package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhangx1n/xim/state"
)

func init() {
	rootCmd.AddCommand(stateCmd)
}

var stateCmd = &cobra.Command{
	Use: "state",
	Run: StateHandle,
}

func StateHandle(cmd *cobra.Command, args []string) {
	state.RunMain(ConfigPath)
}
