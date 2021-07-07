package aws_iam_user_session_response_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/domain_aws/aws_iam_user"
)

// swagger:response getAwsIamUserSessionResponse
type AwsGetIamUserSessionResponseWrapper struct {
	// in: body
	Body AwsGetIamUserSessionResponse
}

type AwsGetIamUserSessionResponse struct {
	Message string
	Data    aws_iam_user.AwsIamUserSession
}

func (responseDto *AwsGetIamUserSessionResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
