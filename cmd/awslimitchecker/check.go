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
		awsService := args[0]
		awsProfile := viper.GetString("awsprofile")
		region := viper.GetString("region")
		console := viper.GetBool("console")

		usage := awslimitchecker.GetLimits(awsService, awsProfile, region)
		if console {
			fmt.Printf("AWS profile: %s | AWS region: %s | service: %s\n", awsProfile, region, awsService)
			for _, u := range usage {
				fmt.Printf("* [%s] %s %g/%g\n",
					u.Service, u.Name, u.UsageValue, u.QuotaValue)
			}
		}
	},
}
