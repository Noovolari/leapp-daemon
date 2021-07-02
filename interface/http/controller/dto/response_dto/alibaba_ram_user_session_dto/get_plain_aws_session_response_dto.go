package plain_alibaba_session_dto

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/domain/session"
)

// swagger:response getAlibabaRamUserSessionResponse
type GetAlibabaRamUserSessionResponseWrapper struct {
  // in: body
  Body GetAlibabaRamUserSessionResponse
}

type GetAlibabaRamUserSessionResponse struct {
  Message string
  Data    session.AlibabaRamUserSession
}

func (responseDto *GetAlibabaRamUserSessionResponse) ToMap() gin.H {
  return gin.H{
    "message": responseDto.Message,
    "data": responseDto.Data,
  }
}
