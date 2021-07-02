package gcp_iam_user_account_oauth_session_dto

import (
	"github.com/gin-gonic/gin"
)

// swagger:response gcpOauthUrlResponse
type GcpOauthUrlResponseWrapper struct {
	// in: body
	Body GcpOauthUrlResponse
}

type GcpOauthUrlResponse struct {
	Message string
	Data    string
}

func (responseDto *GcpOauthUrlResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
