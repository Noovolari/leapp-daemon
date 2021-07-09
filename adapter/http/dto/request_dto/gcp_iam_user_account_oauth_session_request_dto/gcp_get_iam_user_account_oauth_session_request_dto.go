package gcp_iam_user_account_oauth_session_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters getGcpIamUserAccountOauthSession
type GcpGetIamUserAccountOauthSessionRequestDto struct {
  // in: path
  // required: true
  Id string `json:"id" uri:"id" binding:"required"`
}

func (requestDto *GcpGetIamUserAccountOauthSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}