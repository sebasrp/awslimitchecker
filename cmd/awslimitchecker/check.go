package main

import (
	"github.com/sebasrp/awslimitchecker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(check)
}

var check = &cobra.Command{
	Use:   "check",
	Short: "Runc checks on selected services",
	Long:  `Runc checks on selected services. Use all to run all checks`,
	Run: func(cmd *cobra.Command, args []string) {
		awslimitchecker.Gets3Limits(viper.GetString("awsprofile"), viper.GetString("region"))
	},
}
