package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
	@func: WSConnect
	@description:
		handler for endpoint "/api/scheduler/wsconnect"
*/
func (h *Handler) WSConnect(c *gin.Context) {
	// extract client type from url
	clientType, ok := c.GetQuery("type")
	if !ok {
		log.WithFields(
			log.Fields{
				"Client Address": c.Request.Host,
			}).Warn("Failed to extract client type, invalid websocket connection request, abandoned")
		return
	}
	if clientType != model.CLIENT_TYPE_PROVIDER &&
		clientType != model.CLIENT_TYPE_CONSUMER &&
		clientType != model.CLIENT_TYPE_DEPOSITARY &&
		clientType != model.CLIENT_TYPE_FILESTORE {
		log.WithFields(log.Fields{
			"Given Client Type": clientType,
			"Client Address":    c.Request.Host,
		}).Warn("Unknown client type, abandoned")
	}

	// upgrade to websocket connection
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
			"error":          err,
		}).Warn("Failed to upgrade to websocket connection, abandoned")
		return
	}

	ctx := c.Request.Context()

	switch clientType {
	case model.CLIENT_TYPE_CONSUMER:
		// create consumer instance and start to serve it
		consumer, err := h.ConsumerService.CreateConsumer(ctx, ws)
		if err != nil {
			return
		}

		// register receive callbacks based on websocket type
		h.ConsumerService.InitRecvRoute(ctx, consumer)

	case model.CLIENT_TYPE_PROVIDER:
		// create provider instance and start to serve it
		provider, err := h.ProviderService.CreateProvider(ctx, ws)
		if err != nil {
			return
		}

		// register receive callbacks based on websocket type
		h.ProviderService.InitRecvRoute(ctx, provider)

	case model.CLIENT_TYPE_DEPOSITARY:
		// todo

	case model.CLIENT_TYPE_FILESTORE:
		// todo

	default:
		// leave empty
	}
}
