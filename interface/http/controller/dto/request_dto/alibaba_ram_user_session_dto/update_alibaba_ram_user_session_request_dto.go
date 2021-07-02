package plain_alibaba_session_dto

import (
	http_error2 "leapp_daemon/infrastructure/http/http_error"

	"github.com/gin-gonic/gin"
)

// swagger:parameters updateAlibabaRamUserSession
type UpdateAlibabaRamUserSessionUriRequestWrapper struct {
	// plain alibaba session update uri body
	// in:body
	Body UpdateAlibabaRamUserSessionUriRequest
}

// swagger:parameters updateAlibabaRamUserSession
type UpdateAlibabaRamUserSessionRequestWrapper struct {
	// plain alibaba session update uri body
	// in:body
	Body UpdateAlibabaRamUserSessionRequest
}

type UpdateAlibabaRamUserSessionUriRequest struct {
	Id string `uri:"id" binding:"required"`
}

type UpdateAlibabaRamUserSessionRequest struct {
	Name   string `json:"name" binding:"required"`
	Region string `json:"region" binding:"required"`
	//User string `json:"user" binding:"required"`
	AlibabaAccessKeyId     string `json:"alibabaAccessKeyId" binding:"required"`
	AlibabaSecretAccessKey string `json:"alibabaSecretAccessKey" binding:"required"`
	ProfileName            string `json:"profileName"`
}

func (requestDto *UpdateAlibabaRamUserSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *UpdateAlibabaRamUserSessionUriRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}
