/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// natsPublisherCmd represents the natsPublisher command
var natsPublisherCmd = &cobra.Command{
	Use:   "publisher",
	Short: "run nats publisher demo",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return initializeProducer(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("natsPublisher called")
		config.Dump(cmd)
	},
}

func init() {
	natsCmd.AddCommand(natsPublisherCmd)

	natsPublisherCmd.Flags().AddFlagSet(config.GetNatsPublisherFlags())
}

func initializeProducer(cmd *cobra.Command) error {
	config.BindFlags(cmd, viper.GetViper())
	return nil
}
