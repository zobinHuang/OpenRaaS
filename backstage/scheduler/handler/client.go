package handler

import (
	"encoding/json"
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
		return
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
		uuid, ok := c.GetQuery("uuid")
		if !ok {
			log.WithFields(log.Fields{
				"Client Address": c.Request.Host,
			}).Error("Failed to extract uuid, invalid websocket connection request, abandoned")
			return
		}
		// create provider instance and start to serve it
		provider, err := h.ProviderService.CreateProvider(ctx, ws, uuid)
		if err != nil {
			return
		}

		// register receive callbacks based on websocket type
		h.ProviderService.InitRecvRoute(ctx, provider)

	case model.CLIENT_TYPE_DEPOSITARY:
		// not need now

	case model.CLIENT_TYPE_FILESTORE:
		// not need now

	default:
		// leave empty
	}
}

// NodeOnline register node to rds
func (h *Handler) NodeOnline(c *gin.Context) {
	nodeType, ok := c.GetQuery("type")
	if !ok {
		log.WithFields(
			log.Fields{
				"Node Address": c.Request.Host,
			}).Warn("Failed to extract node type, invalid websocket connection request, abandoned")
		return
	}
	if nodeType != model.CLIENT_TYPE_PROVIDER &&
		nodeType != model.CLIENT_TYPE_DEPOSITARY &&
		nodeType != model.CLIENT_TYPE_FILESTORE {
		log.WithFields(log.Fields{
			"Given Node Type": nodeType,
			"Node Address":    c.Request.Host,
		}).Warn("Unknown Node type, abandoned")
		return
	}

	ctx := c.Request.Context()

	switch nodeType {
	case model.CLIENT_TYPE_PROVIDER:
		var headerData model.ProviderCore
		if err := c.BindJSON(&headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error": err.Error(),
				}).Warn("Failed to extract provider json data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.ProviderService.CreateProviderInRDS(ctx, &headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to register provider to rds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case model.CLIENT_TYPE_DEPOSITARY:
		var headerData model.DepositoryCore
		if err := c.BindJSON(&headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error": err.Error(),
				}).Warn("Failed to extract depository json data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.DepositoryService.CreateDepositoryInRDS(ctx, &headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to register depository to rds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case model.CLIENT_TYPE_FILESTORE:
		var headerData model.FileStoreCore
		if err := c.BindJSON(&headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error": err.Error(),
				}).Warn("Failed to extract FileStore json data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.FileStoreService.CreateFileStoreInRDS(ctx, &headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to register FileStore to rds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	default:
		// leave empty
	}
}

// ApplicationOnline register application to rds
func (h *Handler) ApplicationOnline(c *gin.Context) {
	ctx := c.Request.Context()
	var headerData model.StreamApplication
	if err := c.BindJSON(&headerData); err != nil {
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Warn("Failed to extract stream application json data")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.ApplicationService.CreateStreamApplication(ctx, &headerData); err != nil {
		log.WithFields(
			log.Fields{
				"error":      err.Error(),
				"headerData": headerData,
			}).Warn("Failed to register stream application to rds")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if headerData.DepositoryList == "" {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
			"header data":    headerData,
		}).Warn("Failed to get depository list, abandoned")
		c.JSON(http.StatusBadRequest, gin.H{"error": "request reasonable depository list"})
		return
	}
	var depositaryList []string
	err := json.Unmarshal([]byte(headerData.DepositoryList), &depositaryList)
	if err != nil {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
			"error":          err,
			"header data":    headerData,
		}).Warn("Failed to unmarshal DepositoryList, abandoned")
		c.JSON(http.StatusBadRequest, gin.H{"error": "request reasonable depository list"})
		return
	}

	if headerData.FileStoreList == "" {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
			"header data":    headerData,
		}).Warn("Failed to get FileStore list, abandoned")
		c.JSON(http.StatusBadRequest, gin.H{"error": "request reasonable FileStore list"})
		return
	}
	var fileStoreList []string
	err = json.Unmarshal([]byte(headerData.FileStoreList), &fileStoreList)
	if err != nil {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
			"error":          err,
			"header data":    headerData,
		}).Warn("Failed to unmarshal FileStoreList, abandoned")
		c.JSON(http.StatusBadRequest, gin.H{"error": "request reasonable FileStore list"})
		return
	}
}

func (h *Handler) Clear(c *gin.Context) {
	h.ConsumerService.Clear()
}
