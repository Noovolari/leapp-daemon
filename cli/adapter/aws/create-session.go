package aws

import (
  "fmt"
  "leapp_daemon/client"
  "leapp_daemon/client/aws_iam_user_session"
  "leapp_daemon/models"
)


// CreateAwsIamUserSession Creates a new Plain IAM User session in Leapp
func CreateAwsIamUserSession(Name string, AccessKey string, SecretAccessKey string, Region string) error {
  req := models.AwsCreateIamUserSessionRequest{
    AwsAccessKeyID: AccessKey,
    AwsSecretKey:   SecretAccessKey,
    Region:         &Region,
    SessionName:    &Name,
  }

  params := aws_iam_user_session.NewCreateAwsIamUserSessionParams()
  params.Body = &req
  resp, err := client.Default.AwsIamUserSession.CreateAwsIamUserSession(params)
  fmt.Println(resp)
  fmt.Println(err)
  if err != nil {
    return err
  }
  return resp

  //data := map[string]string{
  //  "name": Name,
  //  "awsAccessKeyId": AccessKey,
  //  "awsSecretAccessKey": SecretAccessKey,
  //  "region": Region,
  //}
  //jsonValue, _ := json.Marshal(data)
  //data.Set("name",  Name)
  //data.Set("awsAccessKeyId", AccessKey)
  //data.Set("awsSecretAccessKey", SecretAccessKey)
  //data.Set("region", Region)
  // body := strings.NewReader(data.Encode())

  ///res, err := http.Post(BaseUrl+plainSession, "application/json", bytes.NewBuffer(jsonValue))
  //if err != nil {
    //return err
  //}
  //fmt.Println(res)
  //return nil
}


