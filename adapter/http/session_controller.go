package http

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/adapter/http/dto/request_dto/confirm_mfa_token_request_dto"
	"leapp_daemon/adapter/http/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

//TODO: Blast this struct and refactor in accord with the server architecture!
func (controller *EngineController) ListSession(context *gin.Context) {
	logging.SetContext(context)

	/*requestDto := request_dto2.ListSessionRequestDto{}
	  err := (&requestDto).Build(context)
	  if err != nil {
	  	_ = context.Error(err)
	  	return
	  }

	  listType := requestDto.Type
	  query := requestDto.Query*/

	sessionList, err := use_case.ListAllSessions(
		controller.Providers.GetGcpIamUserAccountOauthSessionFacade(),
		controller.Providers.GetAwsIamUserSessionFacade())
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: sessionList}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

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

func (controller *EngineController) ListNamedProfiles(context *gin.Context) {
	logging.SetContext(context)

	namedProfiles, err := use_case.ListAllNamedProfiles(controller.Providers.GetNamedProfilesFacade())
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: namedProfiles}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
