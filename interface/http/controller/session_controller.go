package controller

import (
	"github.com/gin-gonic/gin"
	logging2 "leapp_daemon/infrastructure/logging"
	confirm_mfa_token_dto2 "leapp_daemon/interface/http/controller/dto/request_dto/confirm_mfa_token_dto"
	response_dto2 "leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (controller *EngineController) ListSession(context *gin.Context) {
	logging2.SetContext(context)

	/*requestDto := request_dto2.ListSessionRequestDto{}
	  err := (&requestDto).Build(context)
	  if err != nil {
	  	_ = context.Error(err)
	  	return
	  }

	  listType := requestDto.Type
	  query := requestDto.Query*/

	sessionList, err := use_case.ListAllSessions(controller.Providers.GetGcpIamUserAccountOauthSessionFacade())
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto2.MessageAndDataResponseDto{Message: "success", Data: sessionList}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) ConfirmMfaToken(context *gin.Context) {
	logging2.SetContext(context)

	requestDto := confirm_mfa_token_dto2.MfaTokenConfirmRequestDto{}
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

	responseDto := response_dto2.MessageAndDataResponseDto{Message: "success", Data: requestDto.SessionId}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) ListNamedProfiles(context *gin.Context) {
	logging2.SetContext(context)

	namedProfiles, err := use_case.ListAllNamedProfiles(controller.Providers.GetNamedProfilesFacade())
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto2.MessageAndDataResponseDto{Message: "success", Data: namedProfiles}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
