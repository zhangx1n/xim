package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

var rootCmd = &cobra.Command{
	Use:   "plato",
	Short: "一个支持百万 QPS 的 IM 系统",
	Run:   Plato,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Plato(cmd *cobra.Command, args []string) {

}

func initConfig() {

}
