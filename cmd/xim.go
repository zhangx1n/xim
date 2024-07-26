package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigPath string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&ConfigPath,
		"config",
		"./xim.yaml",
		"config file (default is ./xim.yaml)")
}

var rootCmd = &cobra.Command{
	Use:   "xim",
	Short: "一个支持百万 QPS 的 IM 系统",
	Run:   xim,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func xim(cmd *cobra.Command, args []string) {

}

func initConfig() {

}
