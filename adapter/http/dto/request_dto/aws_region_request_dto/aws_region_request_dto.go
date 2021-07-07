package aws_region_request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

type AwsRegionRequestDto struct {
	Region    string `json:"region"`
	SessionId string `json:"session_id"`
}

func (requestDto *AwsRegionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
