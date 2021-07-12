package plain_alibaba_session_dto

import (
	"leapp_daemon/infrastructure/http/http_error"

	"github.com/gin-gonic/gin"
)

type DeleteAlibabaRamUserSessionRequest struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *DeleteAlibabaRamUserSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
