package trusted_alibaba_session_dto

import (
	"leapp_daemon/infrastructure/http/http_error"

	"github.com/gin-gonic/gin"
)

// swagger:parameters deleteAlibabaRamRoleChainedSession
type DeleteAlibabaRamRoleChainedSessionRequestDto struct {
	// the id of the trusted alibaba session
	// in: path
	// required: true
	Id string `uri:"id" binding:"required"`
}

func (requestDto *DeleteAlibabaRamRoleChainedSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
