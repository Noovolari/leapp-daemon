package aws_iam_user_session_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters updateAwsIamUserSession
type AwsUpdateIamUserSessionUriRequestWrapper struct {
	// AWS Iam UserName Session update uri body
	// in:body
	Body AwsEditIamUserSessionUriRequest
}

// swagger:parameters updateAwsIamUserSession
type AwsEditIamUserSessionRequestWrapper struct {
	// AWS Iam UserName Session update uri body
	// in:body
	Body AwsEditIamUserSessionRequest
}

type AwsEditIamUserSessionUriRequest struct {
	Id string `uri:"id" binding:"required"`
}

type AwsEditIamUserSessionRequest struct {
	Name               string `json:"name" binding:"required"`
	AccountNumber      string `json:"accountNumber" binding:"required"`
	Region             string `json:"region" binding:"required"`
	User               string `json:"user" binding:"required"`
	MfaDevice          string `json:"mfaDevice"`
	AwsAccessKeyId     string `json:"awsAccessKeyId" binding:"required"`
	AwsSecretAccessKey string `json:"awsSecretAccessKey" binding:"required"`
	ProfileName        string `json:"profileName"`
}

func (requestDto *AwsEditIamUserSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *AwsEditIamUserSessionUriRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
