package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"

	log "github.com/sirupsen/logrus"
)

/*
	@struct: ProviderService
	@description: service layer
*/
type ProviderService struct {
	InstanceRoomDAL model.InstanceRoomDAL
	ProviderDAL     model.ProviderDAL
	ICEServers      string
}

/*
	@struct: ProviderServiceConfig
	@description: used for config instance of struct ProviderService
*/
type ProviderServiceConfig struct {
	ICEServers      string
	ProviderDAL     model.ProviderDAL
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
	}
}

/*
	@func: CreateProvider
	@description:
		create a new provider instance and start to serve it
*/
func (s *ProviderService) CreateProvider(ctx context.Context, ws *websocket.Conn) (*model.Provider, error) {
	// initialize client instance
	providerID := uuid.Must(uuid.NewV4()).String()
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

		// config provider type
		provider.ProviderType = reqPacketData.ProviderType
		log.WithFields(log.Fields{
			"ProviderID":    provider.ClientID,
			"Provider Type": reqPacketData.ProviderType,
		}).Info("Set provider type")

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

		// construct responses to consumers
		respToConsumers := struct {
			TargetDepository string `json:"target_depository"`
			TargetFilestore  string `json:"target_filestore"`
		}{
			TargetDepository: fmt.Sprintf("%s:%s", reqPacketData.SelectedDepository.HostAddress, reqPacketData.SelectedDepository.Port),
			TargetFilestore:  fmt.Sprintf("%s:%s", reqPacketData.SelectedFilestore.HostAddress, reqPacketData.SelectedFilestore.Port),
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
}
