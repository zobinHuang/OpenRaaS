package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
	@struct: ConsumerService
	@description: service layer
*/
type ConsumerService struct {
	ICEServers          string
	ScheduleServiceCore model.ScheduleServiceCore
	ConsumerDAL         model.ConsumerDAL
	ApplicationDAL      model.ApplicationDAL
	InstanceRoomDAL     model.InstanceRoomDAL
}

/*
	@struct: ConsumerServiceConfig
	@description: used for config instance of struct ConsumerService
*/
type ConsumerServiceConfig struct {
	ICEServers          string
	ScheduleServiceCore model.ScheduleServiceCore
	ConsumerDAL         model.ConsumerDAL
	ApplicationDAL      model.ApplicationDAL
	InstanceRoomDAL     model.InstanceRoomDAL
}

/*
	@func: NewConsumerService
	@description:
		create, config and return an instance of struct ConsumerService
*/
func NewConsumerService(c *ConsumerServiceConfig) model.ConsumerService {
	return &ConsumerService{
		ICEServers:          c.ICEServers,
		ScheduleServiceCore: c.ScheduleServiceCore,
		ConsumerDAL:         c.ConsumerDAL,
		ApplicationDAL:      c.ApplicationDAL,
		InstanceRoomDAL:     c.InstanceRoomDAL,
	}
}

/*
	@func: CreateConsumer
	@description:
		create a new consumer instance and start to serve it
*/
func (s *ConsumerService) CreateConsumer(ctx context.Context, ws *websocket.Conn) (*model.Consumer, error) {
	// initialize client instance
	consumerID := uuid.Must(uuid.NewV4()).String()
	sendCallbackList := map[string]func(model.WSPacket){}
	recvCallbackList := map[string]func(model.WSPacket){}
	newConsumer := &model.Consumer{
		Client: model.Client{
			ClientID:            consumerID,
			WebsocketConnection: ws,
			SendCallbackList:    sendCallbackList,
			RecvCallbackList:    recvCallbackList,
			Done:                make(chan struct{}),
		},
	}

	// add to global client list
	s.ConsumerDAL.CreateConsumer(ctx, newConsumer)

	// start to serve it
	go func(consumer *model.Consumer) {
		// listen loop
		consumer.Listen()

		// close websocket connection after Listen() finished
		consumer.Close()
		log.WithFields(log.Fields{
			"ClientID": consumer.ClientID,
		}).Info("Close websocket connection")

		// remove from global list after Listen() finished
		s.ConsumerDAL.DeleteConsumer(ctx, consumer.ClientID)
	}(newConsumer)

	log.WithFields(log.Fields{
		"ClientID": consumerID,
	}).Info("Start to serve for client")

	return newConsumer, nil
}

/*
	@func: InitRecvRoute
	@description:
		initialize receiving callback for consumer instance
*/
func (s *ConsumerService) InitRecvRoute(ctx context.Context, consumer *model.Consumer) {
	/*
		@callback: init_consumer_type
		@description:
			config consumer type
	*/
	consumer.Receive("init_consumer_metadata", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			ConsumerType string `json:"consumer_type"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "init_consumer_metadata",
				"ConsumerID":       consumer.ClientID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// validate type
		if reqPacketData.ConsumerType != model.CONSUMER_TYPE_STREAM && reqPacketData.ConsumerType != model.CONSUMER_TYPE_TERMINAL {
			log.WithFields(log.Fields{
				"Warn Type":           "Recv Callback Error",
				"Recv Packet Type":    "init_consumer_metadata",
				"ConsumerID":          consumer.ClientID,
				"Given Consumer Type": reqPacketData.ConsumerType,
			}).Warn("Unknown client type")
			return model.EmptyPacket
		}

		// config consumer type
		consumer.ConsumerType = reqPacketData.ConsumerType
		log.WithFields(log.Fields{
			"ConsumerID":    consumer.ClientID,
			"Consumer Type": reqPacketData.ConsumerType,
		}).Info("Set consumer type")

		// Notify stream consumer to start initialization
		if reqPacketData.ConsumerType == model.CONSUMER_TYPE_STREAM {
			respPacketData := struct {
				ICEServers string `json:"iceservers"`
				ClientID   string `json:"client_id"`
			}{
				ICEServers: s.ICEServers,
				ClientID:   consumer.ClientID,
			}

			jsonRespPacketData, err := json.Marshal(respPacketData)
			if err != nil {
				log.WithFields(log.Fields{
					"ConsumerID": consumer.ClientID,
				}).Warn("Failed to marshal websocket data when try to notify webrtc initialization, abandoned")
				return model.EmptyPacket
			}

			log.WithFields(log.Fields{
				"ConsumerID": consumer.ClientID,
				"PacketType": "notify_ice_server",
			}).Info("send notification to stream consumer to start stream application initilization")

			return model.WSPacket{
				PacketType: "notify_ice_server",
				Data:       string(jsonRespPacketData),
			}
		}

		return model.EmptyPacket
	})

	/*
		@callback: keep_consumer_alive
		@description:
			heartbeat
	*/
	consumer.Receive("keep_consumer_alive", func(req model.WSPacket) (resp model.WSPacket) {
		// log.WithFields(log.Fields{
		// 	"ConsumerID": consumer.ClientID,
		// }).Info("Consumer Heartbeat")
		return model.EmptyPacket
	})

	/*
		@callback: select_stream_application
		@description:
			select stream application
	*/
	consumer.Receive("select_stream_application", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			ApplicationID  string `json:"application_id"`
			ScreenHeight   string `json:"screen_height"`
			ScreenWidth    string `json:"screen_width"`
			VCodec         string `json:"vcodec"`
			ApplicationFPS string `json:"application_fps"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "select_stream_application",
				"ConsumerID":       consumer.ClientID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.WSPacket{
				PacketType: "invalid_application_data",
				Data:       err.Error(),
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// obtain stream application from application dal
		streamApplication, err := s.ApplicationDAL.GetStreamApplicationByID(ctx, reqPacketData.ApplicationID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":            "Recv Callback Error",
				"Recv Packet Type":     "select_stream_application",
				"ConsumerID":           consumer.ClientID,
				"Given Application ID": reqPacketData.ApplicationID,
				"error":                err,
			}).Warn("Failed to find stream application with given id, abandoned")
			return model.WSPacket{
				PacketType: "invalid_application_data",
				Data:       err.Error(),
			}
		}

		// create stream application instance
		ScreenHeightInt, _ := strconv.Atoi(reqPacketData.ScreenHeight)
		ScreenWidthInt, _ := strconv.Atoi(reqPacketData.ScreenWidth)
		FPSInt, _ := strconv.Atoi(reqPacketData.ApplicationFPS)
		streamInstance := &model.StreamInstance{
			StreamApplication: streamApplication,
			ScreenHeight:      ScreenHeightInt,
			ScreenWidth:       ScreenWidthInt,
			FPS:               FPSInt,
			VCodec:            reqPacketData.VCodec,
		}

		// schedule
		provider, depositaryList, filestoreList, err := s.ScheduleServiceCore.ScheduleStream(ctx, streamInstance)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":            "Recv Callback Error",
				"Recv Packet Type":     "select_stream_application",
				"Given Application ID": reqPacketData.ApplicationID,
				"ConsumerID":           consumer.ClientID,
				"error":                err,
			}).Warn("Failed to schedule stream application, no schedule results, abandoned")

			respToConsumer := struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			}
			respToConsumerString, _ := json.Marshal(respToConsumer)

			return model.WSPacket{
				PacketType: "state_failed_provider_schedule",
				Data:       string(respToConsumerString),
			}
		}

		// generate instance index
		streamInstance.InstanceID = uuid.Must(uuid.NewV4()).String()

		// generate request websocket to provider to start instance
		reqToProvider := struct {
			StreamInstance model.StreamInstance   `json:"stream_instance"`
			DepositaryList []model.DepositaryCore `json:"depositary_list"`
			FilestoreList  []model.FilestoreCore  `json:"filestore_list"`
		}{
			StreamInstance: *streamInstance, // metadata of application instance
			DepositaryList: depositaryList,  // metadata of depositary nodes
			FilestoreList:  filestoreList,   // metadata of filestore nodes
		}

		// register instance room
		s.ScheduleServiceCore.CreateStreamInstanceRoom(ctx, provider, consumer, streamInstance)
		log.WithFields(log.Fields{
			"First ConsumerID":  consumer.ClientID,
			"ProviderID":        provider.ClientID,
			"Stream InstanceID": streamInstance.InstanceID,
		}).Info("Stream instance room created")

		// notify instance initialization
		reqToProviderString, err := json.Marshal(reqToProvider)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "select_stream_application",
				"ConsumerID":       consumer.ClientID,
				"error":            err,
			}).Warn("Failed to marshal request to provider, abandoned")
			return model.WSPacket{
				PacketType: "state_failed_provider_schedule",
				Data:       err.Error(),
			}
		}

		// send notification to provider
		provider.Send(model.WSPacket{
			PacketType: "start_schedule",
			Data:       string(reqToProviderString),
		}, nil)
		log.WithFields(log.Fields{
			"ProviderID":  provider.ClientID,
			"Packet Type": "start_schedule",
		}).Info("Send start schedule notification to selected provider")

		// construct ws packet to consumer
		respToConsumer := struct {
			ProviderID string `json:"provider_id"`
		}{
			ProviderID: provider.ClientID,
		}
		respToConsumerString, err := json.Marshal(respToConsumer)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "select_stream_application",
				"ConsumerID":       consumer.ClientID,
				"error":            err,
			}).Warn("Failed to marshal response to consumer, abandoned")
			return model.WSPacket{
				PacketType: "state_failed_provider_schedule",
				Data:       err.Error(),
			}
		}

		// notify the consumer that provider has been selected
		return model.WSPacket{
			PacketType: "state_provider_scheduled",
			Data:       string(respToConsumerString),
		}
	})

	/*
		@callback: start_streamming
		@description:
			notification of start streaming
	*/
	consumer.Receive("start_streaming", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			InstanceID string `json:"instance_id"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "start_streaming",
				"ConsumerID":       consumer.ClientID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// get provider by given instance id
		provider, err := s.InstanceRoomDAL.GetProviderByInstanceID(ctx, reqPacketData.InstanceID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":         "Recv Callback Error",
				"Recv Packet Type":  "start_streaming",
				"ConsumerID":        consumer.ClientID,
				"Given Instance ID": reqPacketData.InstanceID,
			}).Warn("Can't find corresponding provider based on given instance id, abandoned")

			return model.WSPacket{
				PacketType: "failed_start_streaming",
				Data:       err.Error(),
			}
		}

		// notify provider to start streaming to current consumer
		reqToProvider := struct {
			StreamInstanceID string `json:"stream_instance_id"`
			ConsumerID       string `json:"consumer_id"`
		}{
			StreamInstanceID: reqPacketData.InstanceID,
			ConsumerID:       consumer.ClientID,
		}
		reqToProviderString, err := json.Marshal(reqToProvider)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "start_streaming",
				"ConsumerID":       consumer.ClientID,
				"error":            err,
			}).Warn("Failed to marshal response to consumer, abandoned")
			return model.WSPacket{
				PacketType: "failed_start_streaming",
				Data:       fmt.Errorf("Server internal error: scheduler failed to send start_streaming to provider").Error(),
			}
		}

		// send to provider
		provider.Send(model.WSPacket{
			PacketType: "start_streaming",
			Data:       string(reqToProviderString),
		}, nil)

		log.WithFields(log.Fields{
			"Consumer ID":        consumer.ClientID,
			"Instance ID":        reqPacketData.InstanceID,
			"Target Provider ID": provider.ClientID,
		}).Info("Receive start streaming request from consumer, nofity provider")

		return model.EmptyPacket
	})

	/*
		@callback: answer_sdp
		@description:
			notification of start streaming
	*/
	consumer.Receive("answer_sdp", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			InstanceID string `json:"instance_id"`
			AnswerSDP  string `json:"answer_sdp"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "answer_sdp",
				"ConsumerID":       consumer.ClientID,
				"InstanceID":       reqPacketData.InstanceID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// get provider
		provider, err := s.InstanceRoomDAL.GetProviderByInstanceID(ctx, reqPacketData.InstanceID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "answer_sdp",
				"ConsumerID":       consumer.ClientID,
				"InstanceID":       reqPacketData.InstanceID,
				"error":            err,
			}).Warn("Failed to obtain provider by given instance id, abandoned")
			return model.EmptyPacket
		}

		// construct websocket packet to provider
		var respToProvider = struct {
			ConsumerID string `json:"consumer_id"`
			AnswerSDP  string `json:"answer_sdp"`
		}{
			ConsumerID: consumer.ClientID,
			AnswerSDP:  reqPacketData.AnswerSDP,
		}
		respToProviderString, err := json.Marshal(respToProvider)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "answer_sdp",
				"Consumer ID":      consumer.ClientID,
				"error":            err,
			}).Warn("Failed to marshal response to provider, abandoned")
			return model.EmptyPacket
		}

		// send to provider
		provider.Send(model.WSPacket{
			PacketType: "answer_sdp",
			Data:       string(respToProviderString),
		}, nil)

		log.WithFields(log.Fields{
			"Provider ID": provider.ClientID,
			"Consumer ID": consumer.ClientID,
		}).Info("Forward answer SDP to provider")

		return model.EmptyPacket
	})
}
