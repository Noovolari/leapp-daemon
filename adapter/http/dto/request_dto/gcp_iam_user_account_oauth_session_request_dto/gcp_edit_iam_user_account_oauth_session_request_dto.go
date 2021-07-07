package gcp_iam_user_account_oauth_session_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters editGcpIamUserAccountOauthSession
type GcpEditIamUserAccountOauthSessionUriRequestWrapper struct {
	// gcp iam user account oauth session edit uri
	// in:body
	Body GcpEditIamUserAccountOauthSessionUriRequest
}

// swagger:parameters editGcpIamUserAccountOauthSession
type GcpEditIamUserAccountOauthSessionRequestWrapper struct {
	// gcp gcp-iam-user-account-oauth-session session edit body
	// in:body
	Body GcpEditIamUserAccountOauthSessionRequest
}

type GcpEditIamUserAccountOauthSessionUriRequest struct {
	// the Id of the session
	//required: true
	Id string `uri:"id" binding:"required"`
}

type GcpEditIamUserAccountOauthSessionRequest struct {
	// the name which will be displayed
	// required: true
	Name string `json:"name" binding:"required"`

	// the name of the gcp project
	// required: true
	ProjectName string `json:"projectName" binding:"required"`
}

func (requestDto *GcpEditIamUserAccountOauthSessionUriRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *GcpEditIamUserAccountOauthSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
