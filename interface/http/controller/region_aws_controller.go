package controller

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/domain/region"
  logging2 "leapp_daemon/infrastructure/logging"
  aws_region_dto2 "leapp_daemon/interfaces/http/controller/dto/request_dto/aws_region_dto"
  response_dto2 "leapp_daemon/interfaces/http/controller/dto/response_dto"
  service2 "leapp_daemon/use_cases/service"
  "net/http"
)

func GetAwsRegionListController(context *gin.Context) {
	logging2.SetContext(context)

	responseDto := response_dto2.MessageAndDataResponseDto{Message: "success", Data: region.GetRegionList()}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func EditAwsRegionController(context *gin.Context) {
	logging2.SetContext(context)

	requestDto := aws_region_dto2.AwsRegionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.EditAwsSessionRegion(requestDto.SessionId, requestDto.Region)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto2.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
