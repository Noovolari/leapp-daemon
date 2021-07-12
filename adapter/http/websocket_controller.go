package http

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/adapter/http/dto/response_dto"
  websocket_interface "leapp_daemon/adapter/websocket"
  "leapp_daemon/infrastructure/logging"
  "leapp_daemon/infrastructure/ticker"
  websocket_infrastructure "leapp_daemon/infrastructure/websocket"
  "leapp_daemon/use_case/notification"
  "net/http"
  "time"
)

func (controller *EngineController) RegisterClient(context *gin.Context) {
	logging.SetContext(context)

  websocketConnection, err := controller.Providers.GetWebsocketUpgrader().Upgrade(context.Writer, context.Request, nil)
  if err != nil {
    _ = context.Error(err)
    return
  }

  websocketConnectionWrapper := websocket_infrastructure.NewWebsocketConnectionWrapper(websocketConnection)
  err = websocketConnectionWrapper.Setup()
  if err != nil {
    _ = context.Error(err)
    return
  }

  t := ticker.NewTickerWrapper(time.NewTicker(websocket_infrastructure.PingPeriod))
  subscriber := websocket_interface.NewWebsocketSubscriber(websocketConnectionWrapper, t)

  defer func() {
    if err != nil {
      subscriber.Connection.Close()
    }
  }()

	controller.Providers.GetNotificationHub().Subscribe(subscriber)

  go subscriber.ReadPump()
  go subscriber.WritePump()
}

func (controller *EngineController) Test(context *gin.Context) {
  logging.SetContext(context)

  controller.Providers.GetNotificationHub().BroadcastMessage(notification.Message{
    MessageType: notification.MfaTokenRequest,
    Data:        "test",
  })

  responseDto := response_dto.MessageResponse{Message: "success"}
  context.JSON(http.StatusOK, responseDto.ToMap())
}
