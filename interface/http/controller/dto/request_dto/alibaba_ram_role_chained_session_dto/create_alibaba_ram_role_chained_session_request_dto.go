package trusted_alibaba_session_dto

import (
	"leapp_daemon/infrastructure/http/http_error"

	"github.com/gin-gonic/gin"
)

// swagger:parameters createAlibabaRamRoleChainedSession
type CreateAlibabaRamRoleChainedSessionParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	Body CreateAlibabaRamRoleChainedSessionRequestDto
}

type CreateAlibabaRamRoleChainedSessionRequestDto struct {
	// the parent session id, can be an alibaba plain or federated session
	// it's generated with an uuid v4
	// required: true
	ParentId string `json:"parentId" binding:"required"` //,uuid4

	// the name which will be displayed
	// required: true
	AccountName string `json:"accountName" binding:"required"`

	// the account number of the alibaba account related to the role
	// required: true
	AccountNumber string `json:"accountNumber" binding:"required,numeric,len=16"`

	// the role name
	// required: true
	RoleName string `json:"roleName" binding:"required"`

	// the region on which the session will be initialized
	Region string `json:"region" binding:"required"`

	ProfileName string `json:"profileName"`
}

func (requestDto *CreateAlibabaRamRoleChainedSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
