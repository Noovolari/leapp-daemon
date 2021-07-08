package http_controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/domain_aws"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http_controller/dto/request_dto/aws_region_request_dto"
	"leapp_daemon/interface/http_controller/dto/request_dto/confirm_mfa_token_request_dto"
	"leapp_daemon/interface/http_controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (controller *EngineController) GetNamedProfiles(context *gin.Context) {
	// swagger:route GET /aws/named-profiles namedProfiles getNamedProfiles
	// Get the AWS Named Profiles List
	//   Responses:
	//     200: AwsNamedProfilesResponse

	logging.SetContext(context)

	actions := controller.Providers.GetNamedProfilesActions()
	namedProfiles := actions.GetNamedProfiles()

	responseDto := response_dto.AwsNamedProfilesResponse{
		Message: "success",
		Data:    namedProfiles,
	}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

//TODO refactor with actions
func (controller *EngineController) GetAwsRegionList(context *gin.Context) {
	logging.SetContext(context)

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: domain_aws.GetRegionList()}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

//TODO refactor with actions
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

//TODO refactor with actions
func (controller *EngineController) ConfirmMfaToken(context *gin.Context) {
	logging.SetContext(context)

	requestDto := confirm_mfa_token_request_dto.MfaTokenConfirmRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.RotateSessionCredentialsWithMfaToken(requestDto.SessionId, requestDto.MfaToken)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: requestDto.SessionId}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
