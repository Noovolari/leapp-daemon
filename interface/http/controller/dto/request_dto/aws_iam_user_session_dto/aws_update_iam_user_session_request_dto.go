package aws_iam_user_session_dto

import (
	"github.com/gin-gonic/gin"
	http_error2 "leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters updateAwsIamUserSession
type AwsUpdateIamUserSessionUriRequestWrapper struct {
	// AWS Iam User Session update uri body
	// in:body
	Body AwsUpdateIamUserSessionUriRequest
}

// swagger:parameters updateAwsIamUserSession
type AwsUpdateIamUserSessionRequestWrapper struct {
	// AWS Iam User Session update uri body
	// in:body
	Body AwsUpdateIamUserSessionRequest
}

type AwsUpdateIamUserSessionUriRequest struct {
	Id string `uri:"id" binding:"required"`
}

type AwsUpdateIamUserSessionRequest struct {
	Name               string `json:"name" binding:"required"`
	AccountNumber      string `json:"accountNumber" binding:"required"`
	Region             string `json:"region" binding:"required"`
	User               string `json:"user" binding:"required"`
	MfaDevice          string `json:"mfaDevice"`
	AwsAccessKeyId     string `json:"awsAccessKeyId" binding:"required"`
	AwsSecretAccessKey string `json:"awsSecretAccessKey" binding:"required"`
	ProfileName        string `json:"profileName"`
}

func (requestDto *AwsUpdateIamUserSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *AwsUpdateIamUserSessionUriRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}
