package response_dto

import "github.com/gin-gonic/gin"

// swagger:response messageResponse
type MessageResponseWrapper struct {
	// in: body
	Body MessageResponse
}

type MessageResponse struct {
	Message string
}

func (responseDto *MessageResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
	}
}
