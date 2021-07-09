package http

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/adapter/http/dto/request_dto/gcp_iam_user_account_oauth_session_request_dto"
  "leapp_daemon/adapter/http/dto/response_dto"
  "leapp_daemon/adapter/http/dto/response_dto/gcp_iam_user_account_oauth_session_response_dto"
  "leapp_daemon/infrastructure/logging"
  "net/http"
)

func (controller *EngineController) GetGcpOauthUrl(context *gin.Context) {
	// swagger:route GET /gcp/iam-user-account-oauth-url gcpIamUserAccountOauthSession getGcpOauthUrl
	// Get the GCP OAuth url
	//   Responses:
	//     200: GcpOauthUrlResponse

	logging.SetContext(context)

	actions := controller.Providers.GetGcpIamUserAccountOauthSessionActions()

	url, err := actions.GetOAuthUrl()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := gcp_iam_user_account_oauth_session_response_dto.GcpOauthUrlResponse{
		Message: "success",
		Data:    url,
	}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) CreateGcpIamUserAccountOauthSession(context *gin.Context) {
	// swagger:route POST /gcp/iam-user-account-oauth-sessions gcpIamUserAccountOauthSession createGcpIamUserAccountOauthSession
	// Create a new GCP Iam UserAccount Oauth Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := gcp_iam_user_account_oauth_session_request_dto.GcpCreateIamUserAccountOauthSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpIamUserAccountOauthSessionActions()

	err = actions.CreateSession(requestDto.Name, requestDto.AccountId, requestDto.ProjectName, requestDto.OauthCode)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetGcpIamUserAccountOauthSession(context *gin.Context) {
	// swagger:route GET /gcp/iam-user-account-oauth-sessions/{id} gcpIamUserAccountOauthSession getGcpIamUserAccountOauthSession
	// Get a GCP Iam UserAccount Oauth Session
	//   Responses:
	//     200: GcpGetIamUserAccountOauthSessionResponse

	logging.SetContext(context)

	requestDto := gcp_iam_user_account_oauth_session_request_dto.GcpGetIamUserAccountOauthSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpIamUserAccountOauthSessionActions()
	gcpIamUserAccountOauthSession, err := actions.GetSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := gcp_iam_user_account_oauth_session_response_dto.GcpGetIamUserAccountOauthSessionResponse{
		Message: "success",
		Data:    gcpIamUserAccountOauthSession,
	}

	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartGcpIamUserAccountOauthSession(context *gin.Context) {
	// swagger:route GET /gcp/iam-user-account-oauth-sessions/{id}/start gcpIamUserAccountOauthSession startGcpIamUserAccountOauthSession
	// Start a GCP Iam UserAccount Oauth Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := gcp_iam_user_account_oauth_session_request_dto.GcpStartIamUserAccountOauthSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpIamUserAccountOauthSessionActions()
	err = actions.StartSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopGcpIamUserAccountOauthSession(context *gin.Context) {
	// swagger:route GET /gcp/iam-user-account-oauth-sessions/{id}/stop gcpIamUserAccountOauthSession stopGcpIamUserAccountOauthSession
	// Stop a GCP Iam UserAccount Oauth Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := gcp_iam_user_account_oauth_session_request_dto.GcpStopIamUserAccountOauthSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpIamUserAccountOauthSessionActions()
	err = actions.StopSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteGcpIamUserAccountOauthSession(context *gin.Context) {
	// swagger:route DELETE /gcp/iam-user-account-oauth-sessions/{id} gcpIamUserAccountOauthSession deleteGcpIamUserAccountOauthSession
	// Delete a GCP Iam UserAccount Oauth Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := gcp_iam_user_account_oauth_session_request_dto.GcpDeleteIamUserAccountOauthSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpIamUserAccountOauthSessionActions()
	err = actions.DeleteSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditGcpIamUserAccountOauthSession(context *gin.Context) {
	// swagger:route PUT /gcp/iam-user-account-oauth-sessions/{id} gcpIamUserAccountOauthSession editGcpIamUserAccountOauthSession
	// Edit a GCP Iam UserAccount Oauth Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestUriDto := gcp_iam_user_account_oauth_session_request_dto.GcpEditIamUserAccountOauthSessionUriRequest{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := gcp_iam_user_account_oauth_session_request_dto.GcpEditIamUserAccountOauthSessionRequest{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpIamUserAccountOauthSessionActions()
	err = actions.EditSession(requestUriDto.Id, requestDto.Name, requestDto.ProjectName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
