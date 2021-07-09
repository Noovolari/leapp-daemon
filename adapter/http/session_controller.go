// Package http Leapp API
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /api/v1
//     Version: 0.0.1
//     License: MIT https://opensource.org/licenses/MIT
//     Contact: John Doe<john.doe@example.com> https://john.doe.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
// swagger:meta
package http

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/adapter/http/dto/request_dto/confirm_mfa_token_request_dto"
  "leapp_daemon/adapter/http/dto/response_dto"
  "leapp_daemon/infrastructure/logging"
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
