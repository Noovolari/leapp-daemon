package controller

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	trusted_alibaba_session_dto "leapp_daemon/interface/http/controller/dto/request_dto/alibaba_ram_role_chained_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:response getAlibabaRamRoleChainedSessionResponse
type getAlibabaRamRoleChainedSessionResponseWrapper struct {
	// in: body
	Body getAlibabaRamRoleChainedSessionResponse
}

type getAlibabaRamRoleChainedSessionResponse struct {
	Message string
	Data    session.AlibabaRamRoleChainedSession
}

func (controller *EngineController) CreateAlibabaRamRoleChainedSessionController(context *gin.Context) {
	// swagger:route POST /session/trusted session-trusted-alibaba createAlibabaRamRoleChainedSession
	// Create a new trusted alibaba session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.CreateAlibabaRamRoleChainedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleChainedSessionActions()

	err = actions.Create(requestDto.ParentId,
		requestDto.AccountName,
		requestDto.AccountNumber,
		requestDto.RoleName,
		requestDto.Region,
		requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetAlibabaRamRoleChainedSessionController(context *gin.Context) {
	// swagger:route GET /session/trusted/{id} session-trusted-alibaba getAlibabaRamRoleChainedSession
	// Get a Trusted AWS Session
	//   Responses:
	//     200: getAlibabaRamRoleChainedSessionResponse

	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.GetAlibabaRamRoleChainedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleChainedSessionActions()

	sess, err := actions.Get(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAlibabaRamRoleChainedSessionController(context *gin.Context) {
	// swagger:route PUT /session/trusted/{id} session-trusted-alibaba editAlibabaRamRoleChainedSession
	// Edit a Trusted AWS Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestUriDto := trusted_alibaba_session_dto.EditAlibabaRamRoleChainedSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := trusted_alibaba_session_dto.EditAlibabaRamRoleChainedSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleChainedSessionActions()

	err = actions.Update(
		requestUriDto.Id,
		requestDto.ParentId,
		requestDto.AccountName,
		requestDto.AccountNumber,
		requestDto.RoleName,
		requestDto.Region,
		requestDto.ProfileName)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteAlibabaRamRoleChainedSessionController(context *gin.Context) {
	// swagger:route DELETE /session/trusted/{id} session-trusted-alibaba deleteAlibabaRamRoleChainedSession
	// Delete a Trusted AWS Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.DeleteAlibabaRamRoleChainedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleChainedSessionActions()

	err = actions.Delete(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartAlibabaRamRoleChainedSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.StartAlibabaRamRoleChainedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleChainedSessionActions()

	err = actions.Start(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopAlibabaRamRoleChainedSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.StopAlibabaRamRoleChainedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAlibabaRamRoleChainedSessionActions()

	err = actions.Stop(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
