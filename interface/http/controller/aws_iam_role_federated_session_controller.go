package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller/dto/request_dto/aws_iam_role_federated_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (controller *EngineController) GetAwsIamRoleFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_iam_role_federated_session_dto.AwsGetIamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := use_case.GetAwsIamRoleFederatedSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) CreateAwsIamRoleFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_iam_role_federated_session_dto.AwsCreateIamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.CreateAwsIamRoleFederatedSession(requestDto.Name, requestDto.AccountNumber, requestDto.RoleName,
		requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl,
		requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAwsIamRoleFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestUriDto := aws_iam_role_federated_session_dto.AwsEditIamRoleFederatedSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := aws_iam_role_federated_session_dto.AwsEditIamRoleFederatedSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.UpdateAwsIamRoleFederatedSession(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.RoleName,
		requestDto.RoleArn,
		requestDto.IdpArn,
		requestDto.Region,
		requestDto.SsoUrl,
		requestDto.ProfileName)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteAwsIamRoleFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_iam_role_federated_session_dto.AwsDeleteIamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.DeleteAwsIamRoleFederatedSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartAwsIamRoleFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_iam_role_federated_session_dto.AwsStartIamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StartAwsIamRoleFederatedSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopAwsIamRoleFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_iam_role_federated_session_dto.AwsStopIamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StopAwsIamRoleFederatedSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
