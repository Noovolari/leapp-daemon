package http

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/aws/aws_iam_role_chained"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/adapter/http/dto/request_dto/aws_iam_role_chained_session_request_dto"
	"leapp_daemon/adapter/http/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

// swagger:response getAwsIamRoleChainedSessionResponse
type getAwsIamRoleChainedSessionResponseWrapper struct {
	// in: body
	Body getAwsIamRoleChainedSessionResponse
}

type getAwsIamRoleChainedSessionResponse struct {
	Message string
	Data    aws_iam_role_chained.AwsIamRoleChainedSession
}

func (controller *EngineController) CreateAwsIamRoleChainedSession(context *gin.Context) {
	// swagger:route POST /aws/iam-role-chained-sessions createAwsIamRoleChainedSession
	// Create a new AWS IAM Role Chained Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := aws_iam_role_chained_session_request_dto.AwsCreateIamRoleChainedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.CreateAwsIamRoleChainedSession(requestDto.ParentId, requestDto.AccountName, requestDto.AccountNumber, requestDto.RoleName, requestDto.Region)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetAwsIamRoleChainedSession(context *gin.Context) {
	// swagger:route GET /aws/iam-role-chained-sessions/{id} awsIamRoleChainedSession getAwsIamRoleChainedSession
	// Get a AWS IAM Role Chained Session
	//   Responses:
	//     200: AwsGetIamRoleChainedSessionResponse

	logging.SetContext(context)

	requestDto := aws_iam_role_chained_session_request_dto.AwsGetIamRoleChainedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := use_case.GetAwsIamRoleChainedSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAwsIamRoleChainedSession(context *gin.Context) {
	// swagger:route PUT /aws/iam-role-chained-sessions/{id} awsIamRoleChainedSession editAwsIamRoleChainedSession
	// Edit a AWS IAM Role Chained Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestUriDto := aws_iam_role_chained_session_request_dto.AwsEditIamRoleChainedSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := aws_iam_role_chained_session_request_dto.AwsEditIamRoleChainedSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.UpdateAwsIamRoleChainedSession(
		requestUriDto.Id,
		requestDto.ParentId,
		requestDto.AccountName,
		requestDto.AccountNumber,
		requestDto.RoleName,
		requestDto.Region)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteAwsIamRoleChainedSession(context *gin.Context) {
	// swagger:route DELETE /aws/iam-role-chained-sessions/{id} awsIamRoleChainedSession deleteAwsIamRoleChainedSession
	// Delete a AWS IAM Role Chained Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := aws_iam_role_chained_session_request_dto.AwsDeleteIamRoleChainedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.DeleteAwsIamRoleChainedSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
