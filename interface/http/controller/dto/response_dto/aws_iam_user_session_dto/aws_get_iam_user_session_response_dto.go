package aws_iam_user_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/session"
)

// swagger:response getAwsIamUserSessionResponse
type AwsGetIamUserSessionResponseWrapper struct {
	// in: body
	Body AwsGetIamUserSessionResponse
}

type AwsGetIamUserSessionResponse struct {
	Message string
	Data    session.AwsIamUserSession
}

func (responseDto *AwsGetIamUserSessionResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
