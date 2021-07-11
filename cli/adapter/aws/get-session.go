package aws

import (
	"fmt"
  "leapp_daemon/client"
  "leapp_daemon/client/aws_iam_user_session"
)

// GetAwsIamUserSession Start a IAM User session in Leapp
func GetAwsIamUserSession(id string) error {
  params := aws_iam_user_session.NewGetAwsIamUserSessionParams()
  params.ID = id
  resp, err := client.Default.AwsIamUserSession.GetAwsIamUserSession(params)
  fmt.Println(resp)
  fmt.Println(err)
  if err != nil {
    return err
  }
  return resp
}

