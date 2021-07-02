package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	aws_iam_user_session_request_dto "leapp_daemon/interface/http/controller/dto/request_dto/aws_iam_user_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	aws_iam_user_session_response_dto "leapp_daemon/interface/http/controller/dto/response_dto/aws_iam_user_session_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (controller *EngineController) CreateAwsIamUserSession(context *gin.Context) {
	// swagger:route POST /aws/iam-user-sessions awsIamUserSession createAwsIamUserSession
	// Create a new AWS IAM User Session
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

	err = actions.Create(requestDto.Name, requestDto.AwsAccessKeyId, requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice, requestDto.Region, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetAwsIamUserSession(context *gin.Context) {
	// swagger:route GET /aws/iam-user-sessions/{id} awsIamUserSession getAwsIamUserSession
	// Get a AWS IAM User Session
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

	sess, err := actions.GetAwsIamUserSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := aws_iam_user_session_response_dto.AwsGetIamUserSessionResponse{
		Message: "success",
		Data:    *sess,
	}

	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) UpdateAwsIamUserSession(context *gin.Context) {
	// swagger:route PUT /aws/iam-user-sessions/{id} awsIamUserSession updateAwsIamUserSession
	// Edit a AWS IAM User Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestUriDto := aws_iam_user_session_request_dto.AwsUpdateIamUserSessionUriRequest{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := aws_iam_user_session_request_dto.AwsUpdateIamUserSessionRequest{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsIamUserSessionActions()

	err = actions.UpdateAwsIamUserSession(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.Region,
		requestDto.User,
		requestDto.AwsAccessKeyId,
		requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice,
		requestDto.ProfileName)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteAwsIamUserSession(context *gin.Context) {
	// swagger:route DELETE /aws/iam-user-sessions/{id} awsIamUserSession deleteAwsIamUserSession
	// Delete a AWS IAM User Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := aws_iam_user_session_request_dto.AwsDeleteIamUserSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = session.NewAwsIamUserSessionsFacade().RemoveSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartAwsIamUserSession(context *gin.Context) {
	// swagger:route POST /aws/iam-user-sessions/{id}/start awsIamUserSession startAwsIamUserSession
	// Start a AWS IAM User Session
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

	err = actions.StartAwsIamUserSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopAwsIamUserSession(context *gin.Context) {
	// swagger:route POST /aws/iam-user-sessions/{id}/stop awsIamUserSession stopAwsIamUserSession
	// Stop a AWS IAM User Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := aws_iam_user_session_request_dto.AwsStopIamUserSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StopAwsIamUserSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
