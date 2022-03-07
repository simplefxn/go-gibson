/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/spf13/cobra"
)

// natsCmd represents the nats command
var natsCmd = &cobra.Command{
	Use:   "nats",
	Short: "run nats demo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nats called")
	},
}

func init() {
	runCmd.AddCommand(natsCmd)

	natsCmd.PersistentFlags().AddFlagSet(config.GetNatsGenericFlags())
	natsCmd.PersistentFlags().AddFlagSet(config.GetMetricsFlags())
}
