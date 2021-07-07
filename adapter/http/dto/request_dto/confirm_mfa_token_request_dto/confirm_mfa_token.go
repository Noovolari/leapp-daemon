package confirm_mfa_token_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

type MfaTokenConfirmRequestDto struct {
	SessionId string `json:"sessionId" binding:"required"`
	MfaToken  string `json:"mfaToken" binding:required`
}

func (requestDto *MfaTokenConfirmRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
