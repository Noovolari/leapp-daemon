package gcp_iam_user_account_oauth_session_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters createGcpIamUserAccountOauthSession
type GcpCreateIamUserAccountOauthSessionRequestWrapper struct {
	// gcp iam user account oauth session create body
	// in:body
	Body GcpCreateIamUserAccountOauthSessionRequest
}

type GcpCreateIamUserAccountOauthSessionRequest struct {
	// the name which will be displayed
	// required: true
	Name string `json:"name" binding:"required"`

	// the account identifier of the gcp account
	// required: true
	AccountId string `json:"accountId" binding:"required"`

	// the name of the gcp project
	// required: true
	ProjectName string `json:"projectName" binding:"required"`

	// the OAuth code to obtain credentials
	// required: true
	OauthCode string `json:"oauthCode" binding:"required"`
}

func (requestDto *GcpCreateIamUserAccountOauthSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
