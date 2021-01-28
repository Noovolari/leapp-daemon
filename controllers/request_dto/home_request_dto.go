package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/error_handling"
)

type HomeRequestDto struct {
	Name string `uri:"Name" binding:"required"`
}

func (requestDto *HomeRequestDto) Build(context *gin.Context) error {
	err := error_handling.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}