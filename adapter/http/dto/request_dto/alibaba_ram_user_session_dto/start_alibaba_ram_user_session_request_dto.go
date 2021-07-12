package plain_alibaba_session_dto

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/infrastructure/http/http_error"
)

type StartAlibabaRamUserSessionRequest struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *StartAlibabaRamUserSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
