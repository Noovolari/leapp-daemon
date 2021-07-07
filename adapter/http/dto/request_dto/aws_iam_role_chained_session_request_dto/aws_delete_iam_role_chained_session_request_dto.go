package aws_iam_role_chained_session_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters deleteAwsIamRoleChainedSession
type AwsDeleteIamRoleChainedSessionRequestDto struct {
	// the id of the aws trusted session
	// in: path
	// required: true
	Id string `uri:"id" binding:"required"`
}

func (requestDto *AwsDeleteIamRoleChainedSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
