package response_dto

import "github.com/gin-gonic/gin"

type MessageAndDataResponseDto struct {
	Message string
	Data    interface{}
}

func (responseDto *MessageAndDataResponseDto) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
