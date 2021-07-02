package gcp_iam_user_account_oauth_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/session"
)

// swagger:response gcpGetIamUserAccountOauthSessionResponse
type GcpGetIamUserAccountOauthSessionResponseWrapper struct {
	// in: body
	Body GcpGetIamUserAccountOauthSessionResponse
}

type GcpGetIamUserAccountOauthSessionResponse struct {
	Message string
	Data    session.GcpIamUserAccountOauthSession
}

func (responseDto *GcpGetIamUserAccountOauthSessionResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
