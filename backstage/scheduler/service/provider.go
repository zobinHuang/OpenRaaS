package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/utils"

	"github.com/gorilla/websocket"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"

	log "github.com/sirupsen/logrus"
)

/*
@struct: ProviderService
@description: service layer
*/
type ProviderService struct {
	InstanceRoomDAL model.InstanceRoomDAL
	ProviderDAL     model.ProviderDAL
	ConsumerDAL     model.ConsumerDAL
	ICEServers      string
}

/*
@struct: ProviderServiceConfig
@description: used for config instance of struct ProviderService
*/
type ProviderServiceConfig struct {
	ICEServers      string
	ProviderDAL     model.ProviderDAL
	ConsumerDAL     model.ConsumerDAL
	InstanceRoomDAL model.InstanceRoomDAL
}

/*
@func: NewProviderService
@description:

	create, config and return an instance of struct ProviderService
*/
func NewProviderService(c *ProviderServiceConfig) model.ProviderService {
	return &ProviderService{
		ICEServers:      c.ICEServers,
		ProviderDAL:     c.ProviderDAL,
		InstanceRoomDAL: c.InstanceRoomDAL,
		ConsumerDAL:     c.ConsumerDAL,
	}
}

// CreateProviderInRDS write Provider info to rds
func (s *ProviderService) CreateProviderInRDS(ctx context.Context, provider *model.ProviderCoreWithInst) error {
	return s.ProviderDAL.CreateProviderInRDS(ctx, provider)
}

func (s *ProviderService) UpdateProviderInRDS(ctx context.Context, provider *model.ProviderCoreWithInst) error {
	return s.ProviderDAL.UpdateProviderInRDSByID(ctx, provider)
}

/*
@func: CreateProvider
@description:

	create a new provider instance and start to serve it
*/
func (s *ProviderService) CreateProvider(ctx context.Context, ws *websocket.Conn, providerID string) (*model.Provider, error) {
	// initialize client instance
	sendCallbackList := map[string]func(model.WSPacket){}
	recvCallbackList := map[string]func(model.WSPacket){}
	newProvider := &model.Provider{
		Client: model.Client{
			ClientID:            providerID,
			WebsocketConnection: ws,
			SendCallbackList:    sendCallbackList,
			RecvCallbackList:    recvCallbackList,
			Done:                make(chan struct{}),
		},
	}

	// add to global provider list
	s.ProviderDAL.CreateProvider(ctx, newProvider)

	// start to serve it
	go func(provider *model.Provider) {
		// listen loop
		provider.Listen()

		// close websocket connection after Listen() finished
		provider.Close()
		log.WithFields(log.Fields{
			"ClientID": provider.ClientID,
		}).Info("Close websocket connection")

		// remove from global list after Listen() finished
		s.ProviderDAL.DeleteProvider(ctx, provider.ClientID)
	}(newProvider)

	log.WithFields(log.Fields{
		"ClientID": providerID,
	}).Info("Start to serve for client")

	return newProvider, nil
}

/*
@func: InitRecvRoute
@description:

	initialize receiving callback for provider instance
*/
func (s *ProviderService) InitRecvRoute(ctx context.Context, provider *model.Provider) {
	/*
		@callback: init_provider_type
		@description:
			config provider type
	*/
	provider.Receive("init_provider_metadata", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			ProviderType string `json:"provider_type"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "init_provider_metadata",
				"ProviderID":       provider.ClientID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// validate type
		if reqPacketData.ProviderType != model.PROVIDER_TYPE_STREAM && reqPacketData.ProviderType != model.PROVIDER_TYPE_TERMINAL {
			log.WithFields(log.Fields{
				"Warn Type":           "Recv Callback Error",
				"Recv Packet Type":    "init_provider_metadata",
				"ProviderID":          provider.ClientID,
				"Given Provider Type": reqPacketData.ProviderType,
			}).Warn("Unknown client type")
			return model.EmptyPacket
		}

		// Notify ice servers list to stream provider
		if reqPacketData.ProviderType == model.PROVIDER_TYPE_STREAM {
			respPacketData := struct {
				ICEServers string `json:"iceservers"`
			}{
				ICEServers: s.ICEServers,
			}

			jsonRespPacketData, err := json.Marshal(respPacketData)
			if err != nil {
				log.WithFields(log.Fields{
					"ProviderID": provider.ClientID,
				}).Info("Failed to marshal websocket data when try to notify webrtc initialization, abandoned")
				return model.EmptyPacket
			}

			return model.WSPacket{
				PacketType: "notify_ice_server",
				Data:       string(jsonRespPacketData),
			}
		}

		return model.EmptyPacket
	})

	/*
		@callback: keep_provider_alive
		@description:
			heartbeat
	*/
	provider.Receive("keep_provider_alive", func(req model.WSPacket) (resp model.WSPacket) {
		// log.WithFields(log.Fields{
		// 	"ConsumerID": consumer.ClientID,
		// }).Debug("Consumer Heartbeat")
		return model.EmptyPacket
	})

	/*
		@callback: state_selected_storage
		@description:
			notification from provider of successfully selecting storage node
	*/
	provider.Receive("state_selected_storage", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			StreamInstanceID   string                       `json:"stream_instance_id"`
			SelectedDepository model.DepositoryCoreWithInst `json:"selected_depository"`
			SelectedFileStore  model.FileStoreCoreWithInst  `json:"selected_filestore"`
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
			"Selected Depository": fmt.Sprintf("%s:%s", reqPacketData.SelectedDepository.IP, reqPacketData.SelectedDepository.Port),
			"Selected FileStore":  fmt.Sprintf("%s:%s", reqPacketData.SelectedFileStore.IP, reqPacketData.SelectedFileStore.Port),
		}).Info("Notification from daemon of successfully selecting storage node")

		go func() {
			streamInstanceRoom, err := s.InstanceRoomDAL.GetInstanceRoomByInstanceID(nil, reqPacketData.StreamInstanceID)
			if err != nil {
				log.WithFields(log.Fields{
					"Warn Type":        "Recv Callback Error",
					"Recv Packet Type": "state_selected_storage",
					"error":            err,
					"StreamInstanceID": reqPacketData.StreamInstanceID,
				}).Warn("Failed to GetInstanceRoomByInstanceID during state_selected_storage")
				return
			}
			streamInstanceRoom.SelectedDepository = &model.Depository{DepositoryCore: model.DepositoryCore{ID: reqPacketData.SelectedDepository.ID, Mem: reqPacketData.SelectedDepository.Mem}}
			for _, value := range reqPacketData.SelectedDepository.InstHistory {
				streamInstanceRoom.SelectedDepositoryBandWidth = value
				break
			}
			streamInstanceRoom.SelectedFileStore = &model.FileStore{FileStoreCore: model.FileStoreCore{ID: reqPacketData.SelectedFileStore.ID, Mem: reqPacketData.SelectedFileStore.Mem}}
			for _, value := range reqPacketData.SelectedFileStore.InstHistory {
				streamInstanceRoom.SelectedFileStoreLatency = value
				break
			}
			log.Infof("reqPacketData: %+v", reqPacketData)
			log.Infof("SelectedDepositoryBandWidth: %s, SelectedFileStoreLatency: %s",
				streamInstanceRoom.SelectedDepositoryBandWidth, streamInstanceRoom.SelectedFileStoreLatency)
		}()

		// construct responses to consumers
		respToConsumers := struct {
			DepositoryAddress    string `json:"depository_address"`
			DepositoryIsPowerful bool   `json:"depository_is_powerful"`
			FileStoreAddress     string `json:"filestore_address"`
			FileStoreIsPowerful  bool   `json:"filestore_is_powerful"`
		}{
			DepositoryAddress:    fmt.Sprintf("%s:%s", reqPacketData.SelectedDepository.IP, reqPacketData.SelectedDepository.Port),
			DepositoryIsPowerful: reqPacketData.SelectedDepository.IsContainFastNetspeed,
			FileStoreAddress:     fmt.Sprintf("%s:%s", reqPacketData.SelectedFileStore.IP, reqPacketData.SelectedFileStore.Port),
			FileStoreIsPowerful:  reqPacketData.SelectedFileStore.IsContainFastNetspeed,
		}
		respToConsumersString, err := json.Marshal(respToConsumers)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_selected_storage",
				"error":            err,
			}).Warn("Failed to marshal response to consumer, abandoned")
			return model.EmptyPacket
		}

		// find consumers based on instance id
		consumerMap, err := s.InstanceRoomDAL.GetConsumerMapByInstanceID(ctx, reqPacketData.StreamInstanceID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_selected_storage",
				"error":            err.Error(),
			}).Warn("Failed to obtained consumer map while receiving selected storage node notification from provider, abandoned")
		}

		// notify every consumer in this instance room
		for consumerID := range consumerMap {
			consumer := consumerMap[consumerID]
			consumer.T4 = time.Now()
			log.Infof("%s, 开始实施业务能力动态组合方案, ID: %s", consumer.T4.Format(utils.TIME_LAYOUT), consumer.ClientID)
			log.Infof("T4 = %s", consumer.T4.Sub(consumer.T3))
			consumerMap[consumerID].Send(model.WSPacket{
				PacketType: "state_selected_storage",
				Data:       string(respToConsumersString),
			}, nil)
		}

		return model.EmptyPacket
	})

	/*
		@callback: state_selected_storage
		@description:
			notification from provider of failed to select storage node
	*/
	provider.Receive("state_failed_select_storage", func(req model.WSPacket) (resp model.WSPacket) {
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

		log.WithFields(log.Fields{
			"Stream Instance ID": reqPacketData.StreamInstanceID,
		}).Info("Notification from daemon of failed to select storage node")

		// construct response to consumer
		respToConsumers := struct {
			Error string `json:"error"`
		}{
			Error: reqPacketData.Error,
		}
		respToConsumersString, err := json.Marshal(respToConsumers)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_failed_select_storage",
				"error":            err,
			}).Warn("Failed to marshal response to consumer, abandoned")
			return model.EmptyPacket
		}

		// find consumers based on instance id
		consumerMap, err := s.InstanceRoomDAL.GetConsumerMapByInstanceID(ctx, reqPacketData.StreamInstanceID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_failed_select_storage",
				"error":            err.Error(),
			}).Warn("Failed to obtained consumer map while receiving failure notification of selecting storage node from provider, abandoned")
		}

		// notify every consumer in this instance room
		for consumerID := range consumerMap {
			consumerMap[consumerID].Send(model.WSPacket{
				PacketType: "state_failed_select_storage",
				Data:       string(respToConsumersString),
			}, nil)
		}

		return model.EmptyPacket
	})

	/*
		@callback: state_run_instance
		@description:
			notification from provider of successfully running instance
	*/
	provider.Receive("state_run_instance", func(req model.WSPacket) (resp model.WSPacket) {
		var reqPacketData struct {
			StreamInstanceID string `json:"stream_instance_id"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_run_instance",
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		log.WithFields(log.Fields{
			"Stream Instance ID": reqPacketData.StreamInstanceID,
		}).Info("Provider notified that the instance is now successfully running")

		// construct websocket packet to consumer
		respToConsumers := struct {
			StreamInstanceID string `json:"stream_instance_id"`
		}{
			StreamInstanceID: reqPacketData.StreamInstanceID,
		}
		respToConsumersString, err := json.Marshal(respToConsumers)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_run_instance",
				"error":            err,
			}).Warn("Failed to marshal response to consumer, abandoned")
			return model.EmptyPacket
		}

		// find consumers based on instance id
		consumerMap, err := s.InstanceRoomDAL.GetConsumerMapByInstanceID(ctx, reqPacketData.StreamInstanceID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_run_instance",
				"error":            err.Error(),
			}).Warn("Failed to obtained consumer map while receiving success notification of running instance from provider, abandoned")
		}

		// notify every consumer in this instance room
		for consumerID := range consumerMap {
			consumer := consumerMap[consumerID]
			consumer.T5 = time.Now()
			log.Infof("%s, 完成业务能力动态组合, 开始进行服务推流, ID: %s", consumer.T5.Format(utils.TIME_LAYOUT), consumer.ClientID)
			log.Infof("T5 = %s", consumer.T5.Sub(consumer.T4))
			log.Infof("%s, 服务个性化定制时长 = %s", consumer.T5.Format(utils.TIME_LAYOUT), consumer.T5.Sub(consumer.T0))
			log.Infof("=================")
			log.Infof("客户 ID: %s", consumer.ClientID)
			log.Infof("|T1 = %s | T2 = %s | T3 = %s|", consumer.T1.Sub(consumer.T0), consumer.T2.Sub(consumer.T1), consumer.T3.Sub(consumer.T2))
			log.Infof("|T4 = %s | T5 = %s |", consumer.T4.Sub(consumer.T3), consumer.T5.Sub(consumer.T4))
			log.Infof("调配响应生成时间 = T3 = %s", consumer.T3.Sub(consumer.T2))
			log.Infof("服务个性化定制时长 = T1 + T2 + T3 + T4 + T5 = %s", consumer.T5.Sub(consumer.T0))
			log.Infof("=================")
			consumerMap[consumerID].Send(model.WSPacket{
				PacketType: "state_run_instance",
				Data:       string(respToConsumersString),
			}, nil)
		}

		return model.EmptyPacket
	})

	/*
		@callback: state_failed_run_instance
		@description:
			notification from provider of failed to run instance
	*/
	provider.Receive("state_failed_run_instance", func(req model.WSPacket) (resp model.WSPacket) {
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
			"Stream Instance ID": reqPacketData.StreamInstanceID,
		}).Info("Provider notified that the instance is failed to run")

		// construct websocket packet to consumer
		respToConsumers := struct {
			Error            string `json:"error"`
			StreamInstanceID string `json:"stream_instance_id"`
		}{
			Error:            reqPacketData.Error,
			StreamInstanceID: reqPacketData.StreamInstanceID,
		}
		respToConsumersString, err := json.Marshal(respToConsumers)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_failed_run_instance",
				"error":            err,
			}).Warn("Failed to marshal response to consumer, abandoned")
			return model.EmptyPacket
		}

		// find consumers based on instance id
		consumerMap, err := s.InstanceRoomDAL.GetConsumerMapByInstanceID(ctx, reqPacketData.StreamInstanceID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "state_failed_run_instance",
				"error":            err.Error(),
			}).Warn("Failed to obtained consumer map while receiving failure notification of running instance from provider, abandoned")
		}

		// notify every consumer in this instance room
		for consumerID := range consumerMap {
			consumerMap[consumerID].Send(model.WSPacket{
				PacketType: "state_failed_run_instance",
				Data:       string(respToConsumersString),
			}, nil)
		}

		return model.EmptyPacket
	})

	/*
		@callback: failed_start_streaming
		@description:
			notification from provider of failed to start streaming
	*/
	provider.Receive("failed_start_streaming", func(req model.WSPacket) (resp model.WSPacket) {
		// define request format
		var reqPacketData struct {
			InstanceID string `json:"instance_id"`
			ConsumerID string `json:"consumer_id"`
			Error      string `json:"error"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "failed_start_streaming",
				"ProviderID":       provider.ClientID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// fetch consumer based on given id
		consumer, err := s.ConsumerDAL.GetConsumerByID(ctx, reqPacketData.ConsumerID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":         "Recv Callback Error",
				"Recv Packet Type":  "failed_start_streaming",
				"Given Consumer ID": reqPacketData.ConsumerID,
			}).Warn("Failed to obtain consumer by given consumer id, abandoned")
			return model.EmptyPacket
		}

		// send to consumer
		consumer.Send(model.WSPacket{
			PacketType: "failed_start_streaming",
			Data:       reqPacketData.Error,
		}, nil)

		return model.EmptyPacket
	})

	/*
		@callback: offer_sdp
		@description:
			receive offer SDP, forward to corresponding consumers
	*/
	provider.Receive("offer_sdp", func(req model.WSPacket) (resp model.WSPacket) {
		var reqPacketData struct {
			InstanceID string `json:"instance_id"`
			ConsumerID string `json:"consumer_id"`
			OfferSDP   string `json:"offer_sdp"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "offer_sdp",
				"ProviderID":       provider.ClientID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// obtain consumer by given consumer id
		consumer, err := s.ConsumerDAL.GetConsumerByID(ctx, reqPacketData.ConsumerID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":         "Recv Callback Error",
				"Recv Packet Type":  "offer_sdp",
				"ProviderID":        provider.ClientID,
				"Given Consumer ID": reqPacketData.ConsumerID,
				"error":             err,
			}).Warn("Failed to obtain consumer by given consumer id")
			return model.EmptyPacket
		}

		// construct websocket packet to consumer
		var respPacketData = &struct {
			InstanceID string `json:"instance_id"`
			ConsumerID string `json:"consumer_id"`
			OfferSDP   string `json:"offer_sdp"`
		}{
			InstanceID: reqPacketData.InstanceID,
			ConsumerID: reqPacketData.ConsumerID,
			OfferSDP:   reqPacketData.OfferSDP,
		}

		// marshal websocket packet
		jsonRespPacketData, err := json.Marshal(respPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Provider ID": provider.ClientID,
				"Consumer ID": consumer.ClientID,
			}).Info("Failed to marshal websocket data when try to notify offer SDP to consumer, abandoned")
			return model.EmptyPacket
		}

		// send to consumer
		consumer.Send(model.WSPacket{
			PacketType: "offer_sdp",
			Data:       string(jsonRespPacketData),
		}, nil)

		log.WithFields(log.Fields{
			"Provider ID": provider.ClientID,
			"Consumer ID": consumer.ClientID,
		}).Info("Forward offer SDP to consumer")

		return model.EmptyPacket
	})

	/*
		@callback: provider_ice_candidate
		@description:
			receive provider ICE candidate, forward to corresponding consumers
	*/
	provider.Receive("provider_ice_candidate", func(req model.WSPacket) (resp model.WSPacket) {
		var reqPacketData struct {
			InstanceID           string `json:"instance_id"`
			ConsumerID           string `json:"consumer_id"`
			ProviderICECandidate string `json:"provider_ice_candidate"`
		}

		// parse request
		err := json.Unmarshal([]byte(req.Data), &reqPacketData)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "provider_ice_candidate",
				"Provider ID":      provider.ClientID,
				"Consumer ID":      reqPacketData.ConsumerID,
				"error":            err,
			}).Warn("Failed to decode json during receiving, abandoned")
			return model.EmptyPacket
		}

		// obtain consumer
		consumer, err := s.ConsumerDAL.GetConsumerByID(ctx, reqPacketData.ConsumerID)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":         "Recv Callback Error",
				"Recv Packet Type":  "provider_ice_candidate",
				"Provider ID":       provider.ClientID,
				"Given Consumer ID": reqPacketData.ConsumerID,
			}).Warn("%s, abandoned", err.Error())
			return model.EmptyPacket
		}

		// construct websocket packet to consumer
		respToConsumer := &struct {
			ProviderICECandidate string `json:"provider_ice_candidate"`
		}{
			ProviderICECandidate: reqPacketData.ProviderICECandidate,
		}
		respToConsumerString, err := json.Marshal(respToConsumer)
		if err != nil {
			log.WithFields(log.Fields{
				"Warn Type":        "Recv Callback Error",
				"Recv Packet Type": "provider_ice_candidate",
				"Provider ID":      provider.ClientID,
				"Consumer ID":      consumer.ClientID,
				"error":            err,
			}).Warn("Failed to marshal response to provider, abandoned")
			return model.EmptyPacket
		}

		// send to consumer
		consumer.Send(model.WSPacket{
			PacketType: "provider_ice_candidate",
			Data:       string(respToConsumerString),
		}, nil)

		log.WithFields(log.Fields{
			"Provider ID": provider.ClientID,
			"Consumer ID": consumer.ClientID,
		}).Info("Forward ICE candidates of provider to consumer")

		return model.EmptyPacket
	})
}

// ShowEnterInfo show info when register a new provider
func (s *ProviderService) ShowEnterInfo(ctx context.Context, provider *model.ProviderCoreWithInst) {
	log.Infof("%s, 新服务提供节点上线, ID: %s", utils.GetCurrentTime(), provider.ID)
	log.Infof("详细信息：\n", provider.DetailedInfo())
}

func (s *ProviderService) ShowAllInfo(ctx context.Context) {
	providers, err := s.ProviderDAL.GetProviderInRDS(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("ProviderService ShowAllInfo GetProviderInRDS error")
		return
	}
	log.Infof("整合前，服务提供节点信息：%+v", providers)
	log.Infof("整合后，服务提供节点信息：")
	totalProcessor := 0.0
	powerNum := 0
	normalNum := 0
	abnormalNum := 0
	for _, p := range providers {
		if p.IsContainGPU {
			powerNum += 1
		} else {
			normalNum += 1
		}
		totalProcessor += p.Processor
		if p.GetAbnormalHistoryTimes() != 0 {
			abnormalNum += 1
		}
	}
	log.Infof("%s, 节点数量：%d, 总计算资源：%f GF, 高性能节点数量：%d, 低性能节点数量：%d，异常节点数量： %d",
		utils.GetCurrentTime(), len(providers), totalProcessor, powerNum, normalNum, abnormalNum)
	log.Infof("ID 后有 * 表示服务异常节点")
	s.ProviderDAL.ShowInfoFromRDS(providers)
}
