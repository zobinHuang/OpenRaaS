package communicator

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/BrosCloud/provider/streamer/model"

	"github.com/gorilla/websocket"
)

/*
	@func: NewDaemonConnection
	@description:
		store websocket connection to daemon
*/
func (s *WebsocketCommunicator) NewDaemonConnection(ctx context.Context, conn *websocket.Conn) {
	if s.DaemonWSConnection != nil {
		log.Warn("Websocket connection to daemon already exist, replace with new connection")
	}

	// create websocket model
	sendCallbackList := map[string]func(model.WSPacket){}
	recvCallbackList := map[string]func(model.WSPacket){}
	wsConnection := &model.Websocket{
		WebsocketConnection: conn,
		SendCallbackList:    sendCallbackList,
		RecvCallbackList:    recvCallbackList,
	}

	// record connection
	s.DaemonWSConnection = wsConnection

	// start to serve for websocket connection
	go func(wsConnection *model.Websocket) {
		// listen loop
		wsConnection.Listen()

		// close websocket connection after Listen() finished
		wsConnection.Close()

		log.Info("Close websocket connection to daemon")
	}(wsConnection)

	log.Info("Start to serve websocket to daemon")
}

/*
	@func: InitDaemonRecvRoute
	@description:
		initialize receiving callback
*/
func (s *WebsocketCommunicator) InitDaemonRecvRoute(ctx context.Context) {
	s.DaemonWSConnection.Receive("register_provider_metadata", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			Scheme   string `json:"scheme"`
			Hostname string `json:"hostname"`
			Port     string `json:"port"`
			Path     string `json:"path"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "register_provider_metadata",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// connect to schduler
		err = s.ConnectToScheduler(ctx,
			reqPacketData.Scheme,
			reqPacketData.Hostname,
			reqPacketData.Port,
			reqPacketData.Path,
		)
		if err != nil {
			return model.WSPacket{
				PacketType: "state_failed_connect",
				Data:       err.Error(),
			}
		}

		// start keep alive routine to scheduler
		s.KeepSchedulerConnAlive(ctx)

		// register recv callbacks
		s.InitSchedulerRecvRoute(ctx)

		// generate request websocket to scheduler to register provider metadata
		reqToScheduler := struct{}{}
		reqToSchedulerString, err := json.Marshal(reqToScheduler)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "register_provider_metadata",
				"error":            err,
			}).Warn("Failed to marshal provider metadata, abandoned")
			return model.WSPacket{
				PacketType: "invalid_provider_metadata",
				Data:       err.Error(),
			}
		}

		// send metadata
		s.SchedulerWSConnection.Send(model.WSPacket{
			PacketType: "init_provider_metadata",
			Data:       string(reqToSchedulerString),
		}, nil)

		return model.EmptyPacket
	})
}
