package http_controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/domain_aws"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http_controller/dto/request_dto/aws_region_request_dto"
	"leapp_daemon/interface/http_controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (controller *EngineController) GetAwsRegionList(context *gin.Context) {
	logging.SetContext(context)

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: domain_aws.GetRegionList()}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAwsRegion(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_region_request_dto.AwsRegionRequestDto{}
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

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
