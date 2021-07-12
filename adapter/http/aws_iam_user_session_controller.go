package http

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/adapter/http/dto/request_dto/aws_iam_user_session_request_dto"
  "leapp_daemon/adapter/http/dto/response_dto"
  "leapp_daemon/adapter/http/dto/response_dto/aws_iam_user_session_response_dto"
  "leapp_daemon/infrastructure/logging"
  "net/http"
)

func (controller *EngineController) CreateAwsIamUserSession(context *gin.Context) {
	// swagger:route POST /aws/iam-user-sessions awsIamUserSession createAwsIamUserSession
	// Create a new AWS IAM UserName Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := aws_iam_user_session_request_dto.AwsCreateIamUserSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsIamUserSessionActions()
	err = actions.CreateSession(requestDto.SessionName, requestDto.Region, requestDto.AccountNumber, requestDto.UserName,
		requestDto.AwsAccessKeyId, requestDto.AwsSecretKey, requestDto.MfaDevice, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetAwsIamUserSession(context *gin.Context) {
	// swagger:route GET /aws/iam-user-sessions/{id} awsIamUserSession getAwsIamUserSession
	// Get a AWS IAM UserName Session
	//   Responses:
	//     200: AwsGetIamUserSessionResponse

	logging.SetContext(context)

	requestDto := aws_iam_user_session_request_dto.AwsGetIamUserSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsIamUserSessionActions()
	sess, err := actions.GetSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := aws_iam_user_session_response_dto.AwsGetIamUserSessionResponse{
		Message: "success",
		Data:    sess,
	}

	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartAwsIamUserSession(context *gin.Context) {
	// swagger:route POST /aws/iam-user-sessions/{id}/start awsIamUserSession startAwsIamUserSession
	// Start an AWS IAM UserName Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := aws_iam_user_session_request_dto.AwsStartIamUserSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsIamUserSessionActions()
	err = actions.StartSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopAwsIamUserSession(context *gin.Context) {
	// swagger:route POST /aws/iam-user-sessions/{id}/stop awsIamUserSession stopAwsIamUserSession
	// Stop an AWS IAM UserName Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := aws_iam_user_session_request_dto.AwsStopIamUserSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsIamUserSessionActions()
	err = actions.StopSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAwsIamUserSession(context *gin.Context) {
	// swagger:route PUT /aws/iam-user-sessions/{id} awsIamUserSession editAwsIamUserSession
	// Edit a AWS IAM UserName Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestUriDto := aws_iam_user_session_request_dto.AwsEditIamUserSessionUriRequest{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := aws_iam_user_session_request_dto.AwsEditIamUserSessionRequest{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsIamUserSessionActions()
	err = actions.EditSession(requestUriDto.Id, requestDto.Name, requestDto.Region, requestDto.AccountNumber,
		requestDto.User, requestDto.AwsAccessKeyId, requestDto.AwsSecretAccessKey, requestDto.MfaDevice, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteAwsIamUserSession(context *gin.Context) {
	// swagger:route DELETE /aws/iam-user-sessions/{id} awsIamUserSession deleteAwsIamUserSession
	// Delete an AWS IAM UserName Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := aws_iam_user_session_request_dto.AwsDeleteIamUserSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsIamUserSessionActions()
	err = actions.DeleteSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
