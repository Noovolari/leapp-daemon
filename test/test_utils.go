package test

import (
  "leapp_daemon/infrastructure/http/http_error"
  "testing"
)

func ExpectHttpError(t *testing.T, err error, expectedStatusCode int, expectedError string) {
  customError, isCustomError := err.(http_error.CustomError)
  if !isCustomError {
    t.Fatalf("expected CustomError")
  }
  if customError.StatusCode != expectedStatusCode {
    t.Fatalf("unexpected error status code: %v", customError.StatusCode)
  }
  if customError.Error() != expectedError {
    t.Fatalf("unexpected error: %v", customError.Error())
  }
}
