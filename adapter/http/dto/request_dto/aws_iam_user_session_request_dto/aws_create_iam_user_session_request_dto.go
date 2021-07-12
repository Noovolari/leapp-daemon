package aws_iam_user_session_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters createAwsIamUserSession
type AwsCreateIamUserSessionRequestWrapper struct {
	// AWS Iam UserName Session create body
	// in:body
	Body AwsCreateIamUserSessionRequest
}

type AwsCreateIamUserSessionRequest struct {
	// the name which will be displayed
	// required: true
	SessionName string `json:"sessionName" binding:"required"`

	// the account number of the aws account related to the role
	// required: true
	AccountNumber string `json:"accountNumber" binding:"required,numeric,len=12"`

	// the region on which the session will be initialized
	// required: true
	Region string `json:"region" binding:"required"`

	UserName       string `json:"userName" binding:"required"`
	MfaDevice      string `json:"mfaDevice"`
	AwsAccessKeyId string `json:"awsAccessKeyId" binding:"required"`
	AwsSecretKey   string `json:"awsSecretKey" binding:"required"`
	ProfileName    string `json:"profileName"`
}

func (requestDto *AwsCreateIamUserSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
