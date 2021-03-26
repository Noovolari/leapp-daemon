package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type EchoRequestDto struct {
	Text string `uri:"text" binding:"required"`
}

func (requestDto *EchoRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}