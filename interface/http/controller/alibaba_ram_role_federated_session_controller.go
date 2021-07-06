package controller

import (
	"leapp_daemon/infrastructure/logging"
	federated_alibaba_session_dto "leapp_daemon/interface/http/controller/dto/request_dto/alibaba_ram_role_federated_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (controller *EngineController) GetAlibabaRamRoleFederatedSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.GetAlibabaRamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleFederatedSessionActions()

	sess, err := actions.Get(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) CreateAlibabaRamRoleFederatedSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.CreateAlibabaRamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleFederatedSessionActions()

	err = actions.Create(requestDto.Name, requestDto.RoleName,
		requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl,
		requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAlibabaRamRoleFederatedSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestUriDto := federated_alibaba_session_dto.EditAlibabaRamRoleFederatedSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := federated_alibaba_session_dto.EditAlibabaRamRoleFederatedSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleFederatedSessionActions()

	err = actions.Update(
		requestUriDto.Id,
		requestDto.Name,
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

func (controller *EngineController) DeleteAlibabaRamRoleFederatedSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.DeleteAlibabaRamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleFederatedSessionActions()

	err = actions.Delete(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartAlibabaRamRoleFederatedSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.StartAlibabaRamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleFederatedSessionActions()

	err = actions.Start(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopAlibabaRamRoleFederatedSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.StopAlibabaRamRoleFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleFederatedSessionActions()

	err = actions.Stop(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
