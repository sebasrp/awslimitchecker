package main

import (
	"fmt"

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
	Args: func(cmd *cobra.Command, args []string) error {
		var numArgs = len(args)
		if numArgs != 1 {
			return fmt.Errorf("check command requires to specify a single aws service or `all`. %d were provided", numArgs)
		}
		if !awslimitchecker.IsValidAwsService(args[0]) {
			return fmt.Errorf("invalid aws service provided: %s", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		awslimitchecker.GetLimits(args[0], viper.GetString("awsprofile"), viper.GetString("region"))
	},
}
