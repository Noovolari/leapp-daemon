package aws_iam_role_federated_session_dto

import (
	"github.com/gin-gonic/gin"
	http_error2 "leapp_daemon/infrastructure/http/http_error"
)

type AwsCreateIamRoleFederatedSessionRequestDto struct {
	Name          string `json:"name" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	RoleName      string `json:"roleName" binding:"required"`
	RoleArn       string `json:"roleArn" binding:"required"`
	IdpArn        string `json:"idpArn" binding:"required"`
	Region        string `json:"region" binding:"required"`
	SsoUrl        string `json:"ssoUrl" binding:"required"`
	ProfileName   string `json:"profileName"`
}

func (requestDto *AwsCreateIamRoleFederatedSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}
