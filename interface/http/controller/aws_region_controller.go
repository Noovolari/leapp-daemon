package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/region"
	logging2 "leapp_daemon/infrastructure/logging"
	aws_region_dto2 "leapp_daemon/interface/http/controller/dto/request_dto/aws_region_dto"
	response_dto2 "leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (controller *EngineController) GetAwsRegionList(context *gin.Context) {
	logging2.SetContext(context)

	responseDto := response_dto2.MessageAndDataResponseDto{Message: "success", Data: region.GetAwsRegionList()}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAwsRegion(context *gin.Context) {
	logging2.SetContext(context)

	requestDto := aws_region_dto2.AwsRegionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.EditAwsSessionRegion(requestDto.SessionId, requestDto.Region)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto2.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
