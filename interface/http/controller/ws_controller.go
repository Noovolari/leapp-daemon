package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/infrastructure/websocket"
	"log"
	"net/http"
)

func (controller *EngineController) GetWs(context *gin.Context) {
	logging.SetContext(context)
	serveWs(context.Writer, context.Request)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &websocket.Connection{Send: make(chan websocket.Message, 1), Ws: ws}
	s := websocket.Subscription{Conn: c, Room: websocket.DefaultRoom}

	websocket.Hub.Register <- s

	go s.WritePump()
	go s.ReadPump()
}
