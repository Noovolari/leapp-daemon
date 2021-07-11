package aws_iam_user_session_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters startAwsIamUserSession
type AwsStartIamUserSessionRequest struct {
  // in: path
  // required: true
  ID string `json:"id" uri:"id" binding:"required"`
}

func (requestDto *AwsStartIamUserSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
