package plain_alibaba_session_dto

import (
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_user"

	"github.com/gin-gonic/gin"
)

// swagger:response getAlibabaRamUserSessionResponse
type GetAlibabaRamUserSessionResponseWrapper struct {
	// in: body
	Body GetAlibabaRamUserSessionResponse
}

type GetAlibabaRamUserSessionResponse struct {
	Message string
	Data    alibaba_ram_user.AlibabaRamUserSession
}

func (responseDto *GetAlibabaRamUserSessionResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
