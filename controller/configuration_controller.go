package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controller/response_dto"
	"leapp_daemon/logging"
	"leapp_daemon/service"
	"net/http"
)

func CreateConfigurationController(context *gin.Context) {
	logging.SetContext(context)

	err := service.CreateConfiguration()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func ReadConfigurationController(context *gin.Context) {
	logging.SetContext(context)

	configuration, err := service.ReadConfiguration()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: configuration}
	context.JSON(http.StatusOK, responseDto.ToMap())
}