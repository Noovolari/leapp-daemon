package response_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/domain_gcp/named_configuration"
)

// swagger:response gcpNamedConfigurationsResponse
type GcpNamedConfigurationsResponseWrapper struct {
	// in: body
	Body GcpNamedConfigurationsResponse
}

type GcpNamedConfigurationsResponse struct {
	Message string
	Data    []named_configuration.NamedConfiguration
}

func (responseDto *GcpNamedConfigurationsResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
