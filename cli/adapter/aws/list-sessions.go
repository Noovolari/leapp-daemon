package aws

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
)

// ListAwsIamUserSessions List sessions
func ListAwsIamUserSessions() error {
  resp, err := http.Get(baseUrl+"/sessions")
  println(resp)
  println(err)
  if err != nil {
    return err
  }
  //json.NewDecoder(resp.Body).Decode(
  //  map[string]interface{
  //    "AwsIamUserSessions": []session.AwsIamUserSession,
  //    "GcpSessions": []session.GcpIamUserAccountOauthSession,
  //  }
  //  )

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  var result map[string]interface{}
  json.Unmarshal([]byte(bodyBytes), &result)

  fmt.Printf("%s", bodyBytes)
  return nil
}

