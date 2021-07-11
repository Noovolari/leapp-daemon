package aws

import (
	"fmt"
	"net/http"
)

// const baseUrl string = "http://localhost:8080/api/v1"

// StopIAMUserSession Start a IAM User session in Leapp
func StopIAMUserSession(id string) error {
  _, err := http.Post(baseUrl+"/aws/iam-user-sessions/"+id+"/stop", "application/json", nil)
  if err != nil {
    return err
  }
  fmt.Printf("Session %s started\n", id)
  return nil
}

