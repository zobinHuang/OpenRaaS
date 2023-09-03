package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"serverd/model"

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
	@func: WSConnect
	@description:
		handler for endpoint "/api/daemon/wsconnect"
*/
func (h *Handler) WSConnect(c *gin.Context) {
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
	streamer, err := h.StreamerService.CreateStreamer(ctx, ws)
	if err != nil {
		return
	}

	h.StreamerClient = streamer
	h.InitRecvRoute(ctx)

	h.StreamerService.StateScheduler(ctx, streamer)
}

/*
	@func: InitRecvRoute
	@description: initialize receiving callback
*/
func (h *Handler) InitRecvRoute(ctx context.Context) error {
	streamer := h.StreamerClient

	/*
		@callback: add_wine_instance
		@description:
	*/
	streamer.Receive("add_wine_instance", func(req model.WSPacket) (resp model.WSPacket) {
		instanceModel := &model.InstanceModel{}

		// 1. parse request
		err := json.Unmarshal([]byte(req.Data), &instanceModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "add_wine_instance",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned.")
			return model.EmptyPacket
		}

		fmt.Printf("Unmarshaled request details: %v\n", instanceModel)

		// 2. select storage servers and state to streamer
		err = h.SelectFilestore(ctx, instanceModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "add_wine_instance",
				"error":            err,
			}).Warn("Failed to select filestore.")
			h.StreamerService.SendErrorMsg(ctx, streamer, model.ERROR_TYPE_STORAGE, instanceModel.Instanceid)
			return model.EmptyPacket
		}
		err = h.SelectDepository(ctx, instanceModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "add_wine_instance",
				"error":            err,
			}).Warn("Failed to select depository.")
			h.StreamerService.SendErrorMsg(ctx, streamer, model.ERROR_TYPE_STORAGE, instanceModel.Instanceid)
			return model.EmptyPacket
		}
		h.StreamerService.StateSelectedStorage(ctx, streamer, instanceModel)

		// 3. get new vmid
		err = h.InstanceService.NewVMID(ctx, instanceModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "add_wine_instance",
				"error":            err,
			}).Warn("No more spare vmid.")
			// cannot run instance, but already select storage, so use ERROR_TYPE_INSTANCE
			h.StreamerService.SendErrorMsg(ctx, streamer, model.ERROR_TYPE_INSTANCE, instanceModel.Instanceid)
			return model.EmptyPacket
		}

		// 4. mount and fetch from selected servers
		err = h.MountFilestore(ctx, instanceModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "add_wine_instance",
				"error":            err,
			}).Warn("Failed to mount filestory.")
			// cannot run instance, but already select storage, so use ERROR_TYPE_INSTANCE
			h.StreamerService.SendErrorMsg(ctx, streamer, model.ERROR_TYPE_INSTANCE, instanceModel.Instanceid)
			return model.EmptyPacket
		}

		err = h.FetchDepository(ctx, instanceModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "add_wine_instance",
				"error":            err,
			}).Warn("Failed to fetch depository.")
			// cannot run instance, but already select storage, so use ERROR_TYPE_INSTANCE
			h.StreamerService.SendErrorMsg(ctx, streamer, model.ERROR_TYPE_INSTANCE, instanceModel.Instanceid)
			return model.EmptyPacket
		}

		// 5. create app instance and state to streamer
		err = h.CreateInstanceWithModel(ctx, instanceModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "add_wine_instance",
				"error":            err,
			}).Warn("Failed to build a new instance.")
			h.StreamerService.SendErrorMsg(ctx, streamer, model.ERROR_TYPE_INSTANCE, instanceModel.Instanceid)
			return model.EmptyPacket
		}

		h.StreamerService.StateNewInstance(ctx, streamer, instanceModel)

		return model.EmptyPacket
	})

	/*
		@callback: remove_wine_instance
		@description:
	*/
	streamer.Receive("remove_wine_instance", func(req model.WSPacket) (resp model.WSPacket) {
		var deleteModel *model.DeleteInstanceModel

		// 1. parse request
		err := json.Unmarshal([]byte(req.Data), &deleteModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "remove_wine_instance",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned.")
			return model.EmptyPacket
		}

		// 2. delete instance
		err = h.DeleteInstanceWithModel(ctx, deleteModel)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "remove_wine_instance",
				"error":            err,
			}).Warn("Failed to remove instance.")
			h.StreamerService.SendErrorMsg(ctx, streamer, model.ERROR_TYPE_REMOVE, deleteModel.Instanceid)
			return model.EmptyPacket
		}

		// 3. state delete info
		streamer.Send(model.WSPacket{
			PacketType: "state_remove_instance",
			Data:       "",
		}, nil)

		return model.EmptyPacket
	})

	/*
		@callback: keep_streamer_alive
		@description: heartbeat
	*/
	streamer.Receive("keep_streamer_alive", func(req model.WSPacket) (resp model.WSPacket) {
		// log.WithFields(log.Fields{
		// 	"ConsumerID": consumer.ClientID,
		// }).Debug("Consumer Heartbeat")
		return model.EmptyPacket
	})

	return nil
}
