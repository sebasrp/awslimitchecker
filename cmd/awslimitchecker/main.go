package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile    string
	region     string
	awsprofile string

	rootCmd = &cobra.Command{
		Use:   "awslimitchecker",
		Short: "A cli to retrieve the limits and usage of your aws account",
		Long:  "A cli to retrieve the limits and usage of your aws account",
	}
)

func main() {
	Execute()
}

// Executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// fmt.Println("flag")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.awslimitchecker.yaml)")
	rootCmd.PersistentFlags().StringVar(&awsprofile, "awsprofile", "", "aws profile to use (default `default`)")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "region to evaluate (default `us-east-1`)")

	viper.BindPFlag("awsprofile", rootCmd.PersistentFlags().Lookup("awsprofile"))
	viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))
	viper.SetDefault("awsprofile", "default")
	viper.SetDefault("region", "us-east-1")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".awslimitchecker" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".awslimitchecker")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		fmt.Println("Keys retrieved from file: ", viper.AllKeys())
	}
}
