/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-gibson",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-gibson.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	logger.RegisterLog()

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().AddFlagSet(config.GetGenericFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.SetLogLevel(rootCmd)
	if config.GetConfigFileName() != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config.GetConfigFileName())
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".config" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".config")
		viper.SetConfigName("config")
	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Log.Infof("Using config file: %s", viper.ConfigFileUsed())
	} else {
		logger.Log.Error(err)
		//os.Exit(1)
	}

	viper.SetEnvPrefix(config.EnvPrefix)

	viper.AutomaticEnv() // read in environment variables that match

	config.BindFlags(rootCmd, viper.GetViper())
}
