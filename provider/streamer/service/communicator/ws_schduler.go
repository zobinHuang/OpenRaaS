package communicator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/OpenRaaS/provider/streamer/model"
)

/*
@func: ConnectToScheduler
@description:

	connect to scheduler node
*/
func (s *WebsocketCommunicator) ConnectToScheduler(ctx context.Context, scheme string, hostname string, port string, path string) error {
	// construct host name
	completeHostname := fmt.Sprintf("%s:%s", hostname, port)
	uuid := os.Getenv("SERVER_ID")
	schedulerURL := url.URL{
		Scheme:   scheme,
		Host:     completeHostname,
		Path:     path,
		RawQuery: "type=provider&uuid=" + uuid,
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

		// marshall ice server
		var ICEServers []struct {
			Urls string `json:"urls"`
		}
		err = json.Unmarshal([]byte(reqPacketData.ICEServers), &ICEServers)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "notify_ice_server",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// store ice server request
		for _, iceServer := range ICEServers {
			s.SchedulerDAL.AddICEServers(iceServer.Urls)
			log.WithFields(log.Fields{
				"ICE Servers": iceServer.Urls,
			}).Info("Add ICE Servers")
		}

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
			DepositoryList []model.DepositoryCore `json:"depository_list"`
			FilestoreList  []model.FilestoreCore  `json:"filestore_list"`
		}

		fmt.Printf("start schedule vcodec\n")
		fmt.Printf("%s\n", reqPacketData.StreamInstance.VCodec)

		fmt.Printf("%s", req.Data)

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
			VCodec:         reqPacketData.StreamInstance.VCodec,
			FilestoreList:  reqPacketData.FilestoreList,
			DepositoryList: reqPacketData.DepositoryList,
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

	/*
		@callback: start_streaming
		@description:
			notification of starting streaming
	*/
	s.SchedulerWSConnection.Receive("start_streaming", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			StreamInstanceID string `json:"stream_instance_id"`
			ConsumerID       string `json:"consumer_id"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "start_streaming",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.WSPacket{
				PacketType: "failed_start_streaming",
				Data:       fmt.Errorf("Server internal error: provider failed to start streaming").Error(),
			}
		}

		// get instance by given instance id
		streamInstance, err := s.InstanceDAL.GetStreamInstanceByID(ctx, reqPacketData.StreamInstanceID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":         "Recv Callback Error",
				"Recv Packet Type":  "start_streaming",
				"Given Instance ID": reqPacketData.StreamInstanceID,
				"Given Consumer ID": reqPacketData.ConsumerID,
			}).Warn("Failed to obtain instance by given instance id, abandoned")
			return model.WSPacket{
				PacketType: "failed_start_streaming",
				Data:       fmt.Sprintf("%s", err.Error()),
			}
		}

		// obtain corresponding Pump (create if not exist)
		pump, isExisted := s.WebRTCStreamDAL.GetPumpByInstanceID(ctx, reqPacketData.StreamInstanceID)
		if !isExisted {
			// create new webrtc streamer
			pump, err = s.WebRTCStreamDAL.NewPump(ctx, streamInstance)
			if err != nil {
				log.WithFields(log.Fields{
					"Warn Type":        "Recv Callback Error",
					"Recv Packet Type": "start_streaming",
					"Instance ID":      streamInstance.Instanceid,
					"error":            err.Error(),
				}).Warn("Failed to create Pump, abandoned")
				return model.WSPacket{
					PacketType: "failed_start_streaming",
					Data:       fmt.Errorf("Provider failed to create Pump").Error(),
				}
			}

			log.WithFields(log.Fields{
				"Instance ID": streamInstance.Instanceid,
			}).Info("Create WebRTC Streamer for a new instance")

			// create WebRTC video streamer
			err = pump.CreateVideoListener()
			if err != nil {
				log.WithFields(log.Fields{
					"Warn Type":        "Recv Callback Error",
					"Recv Packet Type": "start_streaming",
					"Instance ID":      streamInstance.Instanceid,
					"error":            err.Error(),
				}).Warn("Failed to create WebRTC video listener, abandoned")
				return model.WSPacket{
					PacketType: "failed_start_streaming",
					Data:       fmt.Errorf("Failed to create WebRTC video listener").Error(),
				}
			}

			// create WebRTC audio streamer
			err = pump.CreateAudioListener()
			if err != nil {
				log.WithFields(log.Fields{
					"Warn Type":        "Recv Callback Error",
					"Recv Packet Type": "start_streaming",
					"Instance ID":      streamInstance.Instanceid,
					"error":            err.Error(),
				}).Warn("Failed to create WebRTC audio listener, abandoned")
				return model.WSPacket{
					PacketType: "failed_start_streaming",
					Data:       fmt.Errorf("Failed to create WebRTC audio listener").Error(),
				}
			}

			// create WebRTC input simulator
			err = pump.CreateInputSimulator()
			if err != nil {
				log.WithFields(log.Fields{
					"Warn Type":        "Recv Callback Error",
					"Recv Packet Type": "start_streaming",
					"Instance ID":      streamInstance.Instanceid,
					"error":            err.Error(),
				}).Warn("Failed to create WebRTC input simulator, abandoned")
				return model.WSPacket{
					PacketType: "failed_start_streaming",
					Data:       fmt.Errorf("Failed to create WebRTC input simulator").Error(),
				}
			}

			// start hijacking video and audio stream
			pump.ListenVideoStream()
			pump.ListenAudioStream()
			log.WithFields(log.Fields{
				"Instance ID": streamInstance.Instanceid,
			}).Info("New WebRTC streamer is now hijacking video and audio stream")

			// start discharging
			pump.Discharge()
			log.WithFields(log.Fields{
				"Instance ID": streamInstance.Instanceid,
			}).Info("New WebRTC streamer is now discharging to downstream WebRTC pipes")

			// start profiling
			if model.ENABLE_PUPMP_PROFILING {
				pump.PerSecondProfiling()
			}
		}

		// create WebRTC Pipe for this consumer
		webRTCPipe, err := s.WebRTCStreamDAL.NewWebRTCPipe(ctx, streamInstance, reqPacketData.ConsumerID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "start_streaming",
				"Instance ID":      streamInstance.Instanceid,
				"error":            err.Error(),
			}).Warn("Failed to create WebRTC Pipe, abandoned")
			return model.WSPacket{
				PacketType: "failed_start_streaming",
				Data:       fmt.Errorf("Failed to create WebRTC Pipe").Error(),
			}
		}

		// open WebRTC pipe
		offerSDP, err := webRTCPipe.Open(s.SchedulerDAL.GetICEServers(), streamInstance.VCodec, func(candidate string) {
			/* This function will be invoked while this WebRTC Pipe received ICE Candidate result */

			// construct provider ice candidate notification to scheudler
			// (would be forwarded to consumer)
			var respToScheduler = &struct {
				InstanceID           string `json:"instance_id"`
				ConsumerID           string `json:"consumer_id"`
				ProviderICECandidate string `json:"provider_ice_candidate"`
			}{
				InstanceID:           reqPacketData.StreamInstanceID,
				ConsumerID:           reqPacketData.ConsumerID,
				ProviderICECandidate: candidate,
			}

			// marshal resp
			respToSchedulerString, err := json.Marshal(respToScheduler)
			if err != nil {
				log.WithFields(log.Fields{
					"Instance ID": reqPacketData.StreamInstanceID,
					"Consumer ID": reqPacketData.ConsumerID,
				}).Warn("Failed to marshal ice candidate of provider which would be sent to scheduler, abandoned")
				return
			}

			log.WithFields(log.Fields{
				"Instance ID": reqPacketData.StreamInstanceID,
			}).Info("Obtain ICE candidate, send to schduler to notify corresponding consumer")

			// send to scheduler
			s.SchedulerWSConnection.Send(model.WSPacket{
				PacketType: "provider_ice_candidate",
				Data:       string(respToSchedulerString),
			}, nil)
		})

		// add newly created pipe to WebRTC streamer
		pump.AddWebRTCPipe(webRTCPipe)
		log.WithFields(log.Fields{
			"Instance ID": streamInstance.Instanceid,
			"Consumer ID": reqPacketData.ConsumerID,
		}).Info("Add WebRTC pipe to upstream pump")

		// start to harvest input data from this pipe
		pump.HarvestInput(webRTCPipe)
		log.WithFields(log.Fields{
			"Instance ID": streamInstance.Instanceid,
			"Consumer ID": reqPacketData.ConsumerID,
		}).Info("Upstream pump start to harvest input data from current WebRTC pipe")

		// construct offer SDP to scheudler
		// (would be forwarded to consumer)
		var reqToScheduler = &struct {
			InstanceID string `json:"instance_id"`
			ConsumerID string `json:"consumer_id"`
			OfferSDP   string `json:"offer_sdp"`
		}{
			InstanceID: reqPacketData.StreamInstanceID,
			ConsumerID: reqPacketData.ConsumerID,
			OfferSDP:   offerSDP,
		}

		// marshal resp
		reqToSchedulerString, err := json.Marshal(reqToScheduler)
		if err != nil {
			log.WithFields(log.Fields{
				"Instance ID": reqPacketData.StreamInstanceID,
				"Consumer ID": reqPacketData.ConsumerID,
			}).Warn("Failed to marshal offer SDP from provider which would be sent to scheduler, abandoned")
			return model.WSPacket{
				PacketType: "failed_start_streaming",
				Data:       fmt.Errorf("Failed to create offer SDP").Error(),
			}
		}

		// send offer SDP to scheduler
		return model.WSPacket{
			PacketType: "offer_sdp",
			Data:       string(reqToSchedulerString),
		}
	})

	/*
		@callback: answer_sdp
		@description:
			answer SDP from consumer
	*/
	s.SchedulerWSConnection.Receive("answer_sdp", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			ConsumerID string `json:"consumer_id"`
			AnswerSDP  string `json:"answer_sdp"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "answer_sdp",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// obtain webRTCPipe
		webRTCPipe, err := s.WebRTCStreamDAL.GetWebRTCPipeByConsumerID(ctx, reqPacketData.ConsumerID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":         "Recv Callback Error",
				"Recv Packet Type":  "answer_sdp",
				"Given Consumer ID": reqPacketData.ConsumerID,
			}).Warn("Failed to obtain WebRTC pipe by given consumer id, abandoned")
			return model.EmptyPacket
		}

		// setRemoteSDP
		err = webRTCPipe.SetRemoteSDP(reqPacketData.AnswerSDP)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "answer_sdp",
				"Consumer ID":      reqPacketData.ConsumerID,
				"Instance ID":      webRTCPipe.StreamInstance.Instanceid,
			}).Warn("Failed to config remote SDP, abandoned")
			return model.EmptyPacket
		}

		return model.EmptyPacket
	})

	/*
		@callback: consumer_ice_candidate
		@description:
			receive consumer ICE candidate
	*/
	s.SchedulerWSConnection.Receive("consumer_ice_candidate", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			ConsumerID           string `json:"consumer_id"`
			ConsumerICECandidate string `json:"consumer_ice_candidate"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "consumer_ice_candidate",
				"Consumer ID":      reqPacketData.ConsumerID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// obtain WebRTC pipe
		webRTCPipe, err := s.WebRTCStreamDAL.GetWebRTCPipeByConsumerID(ctx, reqPacketData.ConsumerID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":         "Recv Callback Error",
				"Recv Packet Type":  "answer_sdp",
				"Given Consumer ID": reqPacketData.ConsumerID,
			}).Warn("Failed to obtain WebRTC pipe by given consumer id, abandoned")
			return model.EmptyPacket
		}

		// add ice candidate
		err = webRTCPipe.AddCandidate(reqPacketData.ConsumerICECandidate)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "answer_sdp",
				"Consumer ID":      reqPacketData.ConsumerID,
				"Instance ID":      webRTCPipe.StreamInstance.Instanceid,
			}).Warn("Failed to add ice candidate of consumer, abandoned")
			return model.EmptyPacket
		}

		return model.EmptyPacket
	})
}
