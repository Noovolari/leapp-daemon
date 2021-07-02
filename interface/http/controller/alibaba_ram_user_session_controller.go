package controller

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	alibaba_ram_user_session_request_dto "leapp_daemon/interface/http/controller/dto/request_dto/alibaba_ram_user_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	alibaba_ram_user_session_response_dto "leapp_daemon/interface/http/controller/dto/response_dto/alibaba_ram_user_session_dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (controller *EngineController) CreateAlibabaRamUserSessionController(context *gin.Context) {
	// swagger:route POST /plain/alibaba/session/ alibabaRamUserSession createAlibabaRamUserSession
	// Create a new Plain Alibaba Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := alibaba_ram_user_session_request_dto.CreateAlibabaRamUserSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamUserSessionActions()

	err = actions.Create(requestDto.Name, requestDto.AlibabaAccessKeyId, requestDto.AlibabaSecretAccessKey,
		requestDto.Region, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetAlibabaRamUserSessionController(context *gin.Context) {
	// swagger:route GET /plain/alibaba/session/{id} alibabaRamUserSession getAlibabaRamUserSession
	// Get a Plain Alibaba Session
	//   Responses:
	//     200: GetAlibabaRamUserSessionResponse

	logging.SetContext(context)

	requestDto := alibaba_ram_user_session_request_dto.GetAlibabaRamUserSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamUserSessionActions()

	sess, err := actions.Get(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := alibaba_ram_user_session_response_dto.GetAlibabaRamUserSessionResponse{
		Message: "success",
		Data:    *sess,
	}

	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) UpdateAlibabaRamUserSessionController(context *gin.Context) {
	// swagger:route PUT /plain/alibaba/session/{id} alibabaRamUserSession updateAlibabaRamUserSession
	// Edit a Plain Alibaba Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestUriDto := alibaba_ram_user_session_request_dto.UpdateAlibabaRamUserSessionUriRequest{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := alibaba_ram_user_session_request_dto.UpdateAlibabaRamUserSessionRequest{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamUserSessionActions()

	err = actions.Update(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.Region,
		//requestDto.User,
		requestDto.AlibabaAccessKeyId,
		requestDto.AlibabaSecretAccessKey,
		requestDto.ProfileName)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteAlibabaRamUserSessionController(context *gin.Context) {
	// swagger:route DELETE /plain/alibaba/session/{id} alibabaRamUserSession deleteAlibabaRamUserSession
	// Delete a Plain Alibaba Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := alibaba_ram_user_session_request_dto.DeleteAlibabaRamUserSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = session.GetAlibabaRamUserSessionsFacade().RemoveSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartAlibabaRamUserSessionController(context *gin.Context) {
	// swagger:route POST /plain/alibaba/session/{id}/start alibabaRamUserSession startAlibabaRamUserSession
	// Start a Plain Alibaba Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := alibaba_ram_user_session_request_dto.StartAlibabaRamUserSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamUserSessionActions()

	err = actions.Start(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopAlibabaRamUserSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := alibaba_ram_user_session_request_dto.StopAlibabaRamUserSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamUserSessionActions()

	err = actions.Stop(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
