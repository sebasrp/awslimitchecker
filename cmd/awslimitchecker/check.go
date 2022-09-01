package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

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
		csvFlag := viper.GetBool("csv")

		usage := awslimitchecker.GetLimits(awsService, awsProfile, region)

		if console {
			fmt.Printf("AWS profile: %s | AWS region: %s | service: %s\n", awsProfile, region, awsService)
			for _, u := range usage {
				fmt.Printf("* [%s] %s %g/%g\n",
					u.Service, u.Name, u.UsageValue, u.QuotaValue)
			}
		}

		if csvFlag {
			csvfile, err := os.Create("awslimitchecker.csv")
			if err != nil {
				fmt.Printf("failed creating file: %s", err)
			}
			csvwriter := csv.NewWriter(csvfile)

			_ = csvwriter.Write([]string{"region", "Service", "Name", "usage", "quota"})
			for _, u := range usage {
				row := []string{region, u.Service, u.Name, strconv.FormatFloat(u.UsageValue, 'f', 2, 64), strconv.FormatFloat(u.QuotaValue, 'f', 2, 64)}
				_ = csvwriter.Write(row)
			}

			csvwriter.Flush()

			csvfile.Close()
		}

	},
}
