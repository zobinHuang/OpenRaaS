package communicator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/BrosCloud/provider/streamer/model"
)

/*
	@func: ConnectToScheduler
	@description:
		connect to scheduler node
*/
func (s *WebsocketCommunicator) ConnectToScheduler(ctx context.Context, scheme string, hostname string, port string, path string) error {
	// construct host name
	completeHostname := fmt.Sprintf("%s:%s", hostname, port)
	schedulerURL := url.URL{
		Scheme:   scheme,
		Host:     completeHostname,
		Path:     path,
		RawQuery: "type=provider",
	}

	if s.SchedulerWSConnection != nil {
		log.Warn("Websocket connection to scheduler already exist, try to reconnect to scheduler")
	}

	// connect to scheduler
	conn, _, err := websocket.DefaultDialer.Dial(schedulerURL.String(), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to build websocket connection to scheduler")

		return err
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
	s.SchedulerWSConnection = wsConnection

	// start to serve for websocket connection
	go func(wsConnection *model.Websocket) {
		// listen loop
		wsConnection.Listen()

		// close websocket connection after Listen() finished
		wsConnection.Close()

		log.Info("Close websocket connection to scheduler")
	}(wsConnection)

	log.Info("Start to serve websocket to scheduler")

	return nil
}

/*
	@func: KeepSchedulerConnAlive
	@description:
		keep alive routine
*/
func (s *WebsocketCommunicator) KeepSchedulerConnAlive(ctx context.Context) {
	go func() {
		timeTicker := time.NewTicker(time.Second * 10)
		for {
			<-timeTicker.C
			s.SchedulerWSConnection.Send(model.WSPacket{
				PacketType: "keep_consumer_alive",
				Data:       "",
			}, nil)
		}
		// timeTicker.Stop()
	}()
}

/*
	@func: InitSchedulerRecvRoute
	@description:
		initialize receiving callback
*/
func (s *WebsocketCommunicator) InitSchedulerRecvRoute(ctx context.Context) {
	/*
		@callback: notify_ice_server
		@description:
			notification of ice server from scheduler
	*/
	s.SchedulerWSConnection.Receive("notify_ice_server", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			ICEServers string `json:"iceservers"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "notify_ice_server",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// store ice server request
		s.SchedulerDAL.SetICEServers(reqPacketData.ICEServers)
		log.WithFields(log.Fields{
			"ICE Servers": reqPacketData.ICEServers,
		}).Info("Set ICE Servers")

		return model.EmptyPacket
	})

	/*
		@callback: start_schedule
		@description:
			notification of starting provider-side schedule
	*/
	s.SchedulerWSConnection.Receive("start_schedule", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			StreamInstance model.StreamInstance   `json:"stream_instance"`
			DepositaryList []model.DepositaryCore `json:"depositary_list"`
			FilestoreList  []model.FilestoreCore  `json:"filestore_list"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "start_schedule",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		log.WithFields(log.Fields{
			"Instance ID":    reqPacketData.StreamInstance.InstanceID,
			"Application ID": reqPacketData.StreamInstance.ApplicationID,
		}).Info("Be nofified to start schedule")

		// construct websocket packet to daemon
		reqToDaemonPacket := &model.StreamInstanceDaemonModel{
			AppPath:        reqPacketData.StreamInstance.ApplicationPath,
			AppFile:        reqPacketData.StreamInstance.ApplicationFile,
			AppName:        reqPacketData.StreamInstance.ApplicationName,
			HWKey:          reqPacketData.StreamInstance.HWKey,
			ScreenWidth:    reqPacketData.StreamInstance.ScreenWidth,
			ScreenHeight:   reqPacketData.StreamInstance.ScreenHeight,
			FPS:            reqPacketData.StreamInstance.FPS,
			FilestoreList:  reqPacketData.FilestoreList,
			DepositaryList: reqPacketData.DepositaryList,
			InstanceCore: model.InstanceCore{
				Instanceid: reqPacketData.StreamInstance.InstanceID,
			},
		}

		// notify provider daemon
		reqToDaemonPacketString, err := json.Marshal(reqToDaemonPacket)
		if err != nil {
			log.WithFields(log.Fields{
				"Instance ID": reqPacketData.StreamInstance.InstanceID,
			}).Warn("Failed to marshal websocket data when try to notify daemon to start schedule, abandoned")
			return model.EmptyPacket
		}
		s.DaemonWSConnection.Send(model.WSPacket{
			PacketType: "add_wine_instance",
			Data:       string(reqToDaemonPacketString),
		}, nil)

		return model.EmptyPacket
	})
}
