package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
	func: WebsocketConnect
	description: handler for endpoint "/api/provider_streamer/wsconnect"
*/
func (h *Handler) WebsocketConnect(c *gin.Context) {
	// upgrade to websocket connection
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"Daemon Address": c.Request.Host,
			"error":          err,
		}).Warn("Failed to upgrade daemon to websocket connection, abandoned")
		return
	}

	ctx := c.Request.Context()

	// store websocket to daemon, and start to serve it
	h.WebsocketCommunicator.NewDaemonConnection(ctx, ws)

	// register recv callbacks
	h.WebsocketCommunicator.InitDaemonRecvRoute(ctx)
}
