package http_controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http_controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

//TODO: refactor these functions in accord with the architecture!
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

	//TODO refactor with actions
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
