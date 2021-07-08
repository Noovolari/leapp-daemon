package response_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/domain_aws/named_profile"
)

// swagger:response awsNamedProfilesResponse
type AwsNamedProfilesResponseWrapper struct {
	// in: body
	Body AwsNamedProfilesResponse
}

type AwsNamedProfilesResponse struct {
	Message string
	Data    []named_profile.NamedProfile
}

func (responseDto *AwsNamedProfilesResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
