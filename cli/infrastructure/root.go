package infrastructure

import (
	"fmt"
	"github.com/spf13/cobra"
	"leapp_daemon/cli/infrastructure/aws"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "leapp",
	Short: "Leapp manages and secures credentials in multi cloud environment",
	Long:  `Leapp manages and secures credentials in multi cloud environment`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
  rootCmd.AddCommand(aws.ListAwsIamUserSessions)
	rootCmd.AddCommand(aws.CreateAwsIamUserSession)
  rootCmd.AddCommand(aws.GetAwsIamUserSession)
}
