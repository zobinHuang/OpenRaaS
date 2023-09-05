package communicator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/OpenRaaS/provider/streamer/model"

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
	@func: KeepDaemonConnAlive
	@description:
		keep alive routine
*/
func (s *WebsocketCommunicator) KeepDaemonConnAlive(ctx context.Context) {
	go func() {
		timeTicker := time.NewTicker(time.Second * 10)
		for {
			<-timeTicker.C
			s.DaemonWSConnection.Send(model.WSPacket{
				PacketType: "keep_consumer_alive",
				Data:       "",
			}, nil)
		}
		// timeTicker.Stop()
	}()
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

		log.WithFields(log.Fields{
			"Scheduler Scheme":   reqPacketData.Scheme,
			"Scheduler Hostname": reqPacketData.Hostname,
			"Scheduler Port":     reqPacketData.Port,
			"Scheduler Path":     reqPacketData.Path,
		}).Info("Provider daemon indicates metadata of scheduler")

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

		// generate request to scheduler to register provider metadata
		reqToScheduler := struct {
			ProviderType string `json:"provider_type"`
		}{
			ProviderType: "stream",
		}
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

	/*
		@callback: state_selected_storage
		@description:
			notification of successfully select storage nodes (i.e. both depository and filestore)
	*/
	s.DaemonWSConnection.Receive("state_selected_storage", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			StreamInstanceID   string               `json:"stream_instance_id"`
			SelectedDepository model.DepositaryCore `json:"selected_depository"`
			SelectedFilestore  model.FilestoreCore  `json:"selected_filestore"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_selected_storage",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		log.WithFields(log.Fields{
			"Stream Instance ID":  reqPacketData.StreamInstanceID,
			"Selected Depository": fmt.Sprintf("%s:%s", reqPacketData.SelectedDepository.HostAddress, reqPacketData.SelectedDepository.Port),
			"Selected Filestore":  fmt.Sprintf("%s:%s", reqPacketData.SelectedFilestore.HostAddress, reqPacketData.SelectedFilestore.Port),
		}).Info("Notification from daemon of successfully selecting storage node")

		// construct websocket packet to scheduler
		reqToScheduler := struct {
			StreamInstanceID   string               `json:"stream_instance_id"`
			SelectedDepository model.DepositaryCore `json:"selected_depository"`
			SelectedFilestore  model.FilestoreCore  `json:"selected_filestore"`
		}{
			StreamInstanceID:   reqPacketData.StreamInstanceID,
			SelectedDepository: reqPacketData.SelectedDepository,
			SelectedFilestore:  reqPacketData.SelectedFilestore,
		}

		// notify scheduler
		reqToSchedulerPacketString, err := json.Marshal(reqToScheduler)
		if err != nil {
			log.WithFields(log.Fields{
				"Instance ID": reqPacketData.StreamInstanceID,
			}).Warn("Failed to marshal websocket data when try to notify scheudler the success of finding proper storage node, abandoned")
			return model.EmptyPacket
		}
		s.SchedulerWSConnection.Send(model.WSPacket{
			PacketType: "state_selected_storage",
			Data:       string(reqToSchedulerPacketString),
		}, nil)

		return model.EmptyPacket
	})

	/*
		@callback: state_failed_select_storage
		@description:
			notification of failed to select storage nodes (i.e. both depository and filestore)
	*/
	s.DaemonWSConnection.Receive("state_failed_select_storage", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			StreamInstanceID string `json:"stream_instance_id"`
			Error            string `json:"error"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_failed_select_storage",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// construct websocket packet to scheduler
		reqToScheduler := struct {
			StreamInstanceID string `json:"stream_instance_id"`
			Error            string `json:"error"`
		}{
			StreamInstanceID: reqPacketData.StreamInstanceID,
			Error:            reqPacketData.Error,
		}

		// notify scheduler
		reqToSchedulerPacketString, err := json.Marshal(reqToScheduler)
		if err != nil {
			log.WithFields(log.Fields{
				"Instance ID": reqPacketData.StreamInstanceID,
			}).Warn("Failed to marshal websocket data when try to notify scheudler the failure of finding proper storage node, abandoned")
			return model.EmptyPacket
		}
		s.SchedulerWSConnection.Send(model.WSPacket{
			PacketType: "state_failed_select_storage",
			Data:       string(reqToSchedulerPacketString),
		}, nil)

		return model.EmptyPacket
	})

	/*
		@callback: state_run_instance
		@description:
			notification of successfully prepare instance
	*/
	s.DaemonWSConnection.Receive("state_run_instance", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		// note: should be expecifly allocated in heap!
		var reqPacketData = &model.StreamInstanceDaemonModel{}

		// parse request
		err := json.Unmarshal([]byte(req.Data), reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_run_instance",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		fmt.Printf("%v", reqPacketData)

		// store to instance dal
		s.InstanceDAL.AddNewStreamInstance(ctx, reqPacketData)

		log.WithFields(log.Fields{
			"Instance ID": reqPacketData.Instanceid,
		}).Info("Daemon notified that the instance is now successfully running")

		// construct websocket packet to scheduler
		respToSchduler := struct {
			StreamInstanceID string `json:"stream_instance_id"`
		}{
			StreamInstanceID: reqPacketData.Instanceid,
		}

		// notify scheduler
		respToSchdulerPacketString, err := json.Marshal(respToSchduler)
		if err != nil {
			log.WithFields(log.Fields{
				"Instance ID": reqPacketData.Instanceid,
			}).Warn("Failed to marshal websocket data when try to notify scheudler the success of running instance, abandoned")
			return model.EmptyPacket
		}
		s.SchedulerWSConnection.Send(model.WSPacket{
			PacketType: "state_run_instance",
			Data:       string(respToSchdulerPacketString),
		}, nil)

		return model.EmptyPacket
	})

	/*
		@callback: state_failed_run_instance
		@description:
			notification of failed to run instance
	*/
	s.DaemonWSConnection.Receive("state_failed_run_instance", func(req model.WSPacket) (resp model.WSPacket) {
		var reqPacketData struct {
			Error            string `json:"error"`
			StreamInstanceID string `json:"stream_instance_id"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_failed_run_instance",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		log.WithFields(log.Fields{
			"Instance ID": reqPacketData.StreamInstanceID,
		}).Info("Daemon notified that the instance is failed to run")

		// construct websocket packet to scheduler
		respToSchduler := struct {
			Error            string `json:"error"`
			StreamInstanceID string `json:"stream_instance_id"`
		}{
			Error:            reqPacketData.Error,
			StreamInstanceID: reqPacketData.StreamInstanceID,
		}

		// notify scheduler
		respToSchdulerPacketString, err := json.Marshal(respToSchduler)
		if err != nil {
			log.WithFields(log.Fields{
				"Instance ID": reqPacketData.StreamInstanceID,
			}).Warn("Failed to marshal websocket data when try to notify scheudler the success of running instance, abandoned")
			return model.EmptyPacket
		}
		s.SchedulerWSConnection.Send(model.WSPacket{
			PacketType: "state_failed_run_instance",
			Data:       string(respToSchdulerPacketString),
		}, nil)

		return model.EmptyPacket
	})
}
