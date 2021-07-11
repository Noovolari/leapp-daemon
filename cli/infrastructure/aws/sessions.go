package aws

import (
  "github.com/spf13/cobra"
  "leapp_daemon/cli/adapter/aws"
)

var (
	//Type    string // Type of session: supported strings [plain, chained, federated]
	//ValidTypesValues = [3]string{"plain", "chained", "federated"}

	Name string
	AccessKey string
	SecretAccessKey string
	Region string

	ListAwsIamUserSessions = &cobra.Command{
    Use:   "aws-list-user-sessions",
    Short: "Creates a new Leapp session",
    Long:  `Creates a new Plain/Trusted/Federated Leapp session`,
    Run: func(cmd *cobra.Command, args []string) {

      aws.ListAwsIamUserSessions()
    },
  }

	CreateAwsIamUserSession = &cobra.Command{
		Use:   "aws-create-session",
		Short: "Creates a new Leapp session",
		Long:  `Creates a new Plain/Trusted/Federated Leapp session`,
		Run: func(cmd *cobra.Command, args []string) {

      aws.CreateAwsIamUserSession(Name, AccessKey, SecretAccessKey, Region)
		},
	}

  GetAwsIamUserSession = &cobra.Command{
    Use:   "aws-get-session",
    Short: "List a new Leapp session",
    Long:  `List a new Plain Leapp session`,
    Run: func(cmd *cobra.Command, args []string) {
      _ = aws.GetAwsIamUserSession(args[0])
    },
  }

)

func init() {
  CreateAwsIamUserSession.Flags().StringVarP(&Name, "name", "n", "", "name")
  CreateAwsIamUserSession.MarkFlagRequired("name")

  CreateAwsIamUserSession.Flags().StringVarP(&AccessKey, "access-key", "k", "", "access key")
  CreateAwsIamUserSession.MarkFlagRequired("access-key")

  CreateAwsIamUserSession.Flags().StringVarP(&SecretAccessKey, "secret-access-key", "s", "", "secret access key")
  CreateAwsIamUserSession.MarkFlagRequired("secret-access-key")

  CreateAwsIamUserSession.Flags().StringVarP(&Region, "region", "r", "eu-west-1", "region")


  //NewIamUserSession.Flags().IntVarP(&AccountID, "account-id", "n", 0, "account id number")
  //NewIamUserSession.MarkFlagRequired("account-id")

}



