package main

import (
	"fmt"

	"github.com/nyambati/aws-service-limits-exporter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(iam)
}

var iam = &cobra.Command{
	Use:   "iam",
	Short: "Returns necessary iam policies to retrieve usage/limits",
	Long:  `Returns necessary iam policies to retrieve usage/limits`,
	Run: func(cmd *cobra.Command, args []string) {
		iamPolicies := awslimitchecker.GetIamPolicies()
		fmt.Print("Required IAM permissions to retrieve usage/limits:\n")
		for _, p := range iamPolicies {
			fmt.Printf("* %s\n", p)
		}
	},
}
