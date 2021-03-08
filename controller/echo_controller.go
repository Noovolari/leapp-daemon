package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controller/request_dto"
	"leapp_daemon/controller/response_dto"
	"leapp_daemon/logging"
	"leapp_daemon/service"
	"net/http"
)

func EchoController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := request_dto.EchoRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	serviceResponse, err2 := service.Echo(requestDto.Text)
	if err2 != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: serviceResponse}
	context.JSON(http.StatusOK, responseDto.ToMap())
}