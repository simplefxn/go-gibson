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

// natsSubscriberCmd represents the natsSubscriber command
var natsSubscriberCmd = &cobra.Command{
	Use:   "subscriber",
	Short: "run nats subscriber demo",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return initializeSubscriber(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("natsSubscriber called")
	},
}

func init() {
	natsCmd.AddCommand(natsSubscriberCmd)
	natsPublisherCmd.Flags().AddFlagSet(config.GetNatsSubscriberFlags())
}

func initializeSubscriber(cmd *cobra.Command) error {
	config.BindFlags(cmd, viper.GetViper())
	return nil
}
