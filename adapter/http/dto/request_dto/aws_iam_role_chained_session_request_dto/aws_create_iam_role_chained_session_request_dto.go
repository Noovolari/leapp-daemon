package aws_iam_role_chained_session_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters createAwsIamRoleChainedSession
type AwsCreateIamRoleChainedSessionParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	Body AwsCreateIamRoleChainedSessionRequestDto
}

type AwsCreateIamRoleChainedSessionRequestDto struct {
	// the parent session id, can be an aws plain or federated session
	// it's generated with an uuid v4
	// required: true
	ParentId string `json:"parentId" binding:"required,uuid4"`

	// the name which will be displayed
	// required: true
	AccountName string `json:"accountName" binding:"required"`

	// the account number of the aws account related to the role
	// required: true
	AccountNumber string `json:"accountNumber" binding:"required,numeric,len=12"`

	// the role name
	// required: true
	RoleName string `json:"roleName" binding:"required"`

	// the region on which the session will be initialized
	Region string `json:"region" binding:"awsregion"`
}

func (requestDto *AwsCreateIamRoleChainedSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
