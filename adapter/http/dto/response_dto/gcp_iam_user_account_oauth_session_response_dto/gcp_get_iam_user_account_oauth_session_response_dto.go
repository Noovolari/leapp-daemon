package gcp_iam_user_account_oauth_session_response_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/gcp/gcp_iam_user_account_oauth"
)

// swagger:response gcpGetIamUserAccountOauthSessionResponse
type GcpGetIamUserAccountOauthSessionResponseWrapper struct {
	// in: body
	Body GcpGetIamUserAccountOauthSessionResponse
}

type GcpGetIamUserAccountOauthSessionResponse struct {
	Message string
	Data    gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
}

func (responseDto *GcpGetIamUserAccountOauthSessionResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
