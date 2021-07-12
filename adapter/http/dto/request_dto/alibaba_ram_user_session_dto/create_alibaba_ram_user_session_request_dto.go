package plain_alibaba_session_dto

import (
	"leapp_daemon/infrastructure/http/http_error"

	"github.com/gin-gonic/gin"
)

// swagger:parameters createAlibabaRamUserSession
type CreateAlibabaRamUserSessionRequestWrapper struct {
	// plain alibaba session create body
	// in:body
	Body CreateAlibabaRamUserSessionRequest
}

type CreateAlibabaRamUserSessionRequest struct {
	// the name which will be displayed
	// required: true
	Name string `json:"name" binding:"required"`

	// the region on which the session will be initialized
	// required: true
	Region string `json:"region" binding:"required"`

	/*User                   string `json:"user" binding:"required"`*/
	AlibabaAccessKeyId     string `json:"alibabaAccessKeyId" binding:"required"`
	AlibabaSecretAccessKey string `json:"alibabaSecretAccessKey" binding:"required"`
	ProfileName            string `json:"profileName"`
}

func (requestDto *CreateAlibabaRamUserSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
