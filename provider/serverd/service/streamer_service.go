package service

import (
	"context"
	"encoding/json"
	"fmt"
	"serverd/model"
	"serverd/utils"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

/*
	struct: StreamerService
	description: streamer service layer
*/
type StreamerService struct {
	StreamerDone chan struct{}
}

/*
	struct: StreamerServiceConfig
	description: used for config instance of struct StreamerService
*/
type StreamerServiceConfig struct {
}

/*
	func: NewStreamerService
	description: create, config and return an instance of struct StreamerService
*/
func NewStreamerService(c *StreamerServiceConfig) model.StreamerService {
	ss := &StreamerService{}
	return ss
}

/*
	func: RunStreamerContainer
	description: run a provider streamer docker container on the host
*/
func (c *StreamerService) RunStreamerContainer(ctx context.Context) error {
	fmt.Println("invoked(RunStreamerCountainer)")

	// run container
	var execCmd string
	var params []string

	// Add Exec command
	execCmd = "sh"
	params = append(params, "../run-streamer.sh")

	// Add params
	done := utils.RunShell(execCmd, params)

	if done == nil {
		fmt.Println("Failed to run streamer container.")
		return fmt.Errorf("failed to run streamer container")
	}

	fmt.Println("Succeed to run streamer container.")

	c.StreamerDone = done

	return nil
}

/*
	func: KillStreamerContainer
	description: kill the provider streamer docker container
*/
func (c *StreamerService) KillStreamerContainer(ctx context.Context) error {
	fmt.Println("invoked(KillStreamerContainer)")

	c.StreamerDone <- struct{}{}

	// close container
	var execCmd string
	var params []string

	// Add Exec command
	execCmd = "sh"
	params = append(params, "../stop-streamer.sh")

	utils.RunShellWithReturn(execCmd, params)

	return nil
}

/*
	@func: CreateStreamer
	@description:
		create a new provider streamer instance and start to serve it
*/
func (s *StreamerService) CreateStreamer(ctx context.Context, ws *websocket.Conn) (*model.Streamer, error) {
	// initialize client instance
	streamerID := uuid.Must(uuid.NewV4()).String()
	sendCallbackList := map[string]func(model.WSPacket){}
	recvCallbackList := map[string]func(model.WSPacket){}
	newStreamer := &model.Streamer{
		Client: model.Client{
			ClientID:            streamerID,
			WebsocketConnection: ws,
			SendCallbackList:    sendCallbackList,
			RecvCallbackList:    recvCallbackList,
			Done:                make(chan struct{}),
		},
		StreamerCore: model.StreamerCore{},
	}

	// start to serve it
	go func(streamer *model.Streamer) {
		// listen loop
		streamer.Listen()

		// close websocket connection after Listen() finished
		streamer.Close()
		log.WithFields(log.Fields{
			"ClientID": streamer.ClientID,
		}).Info("Close websocket connection")
	}(newStreamer)

	log.WithFields(log.Fields{
		"ClientID": streamerID,
	}).Info("Start to serve for client")

	return newStreamer, nil
}

/*
	func: StateScheduler
	description: use websocket to send information of the scheduler server
*/
func (c *StreamerService) StateScheduler(ctx context.Context, streamer *model.Streamer) error {
	// load configure data
	config := &model.Config{}
	config.LoadConfigFile()

	// define request format
	reqData := struct {
		SchedulerHost   string `json:"hostname"`
		SchedulerPath   string `json:"path"`
		SchedulerPort   string `json:"port"`
		SchedulerScheme string `json:"scheme"`
	}{
		SchedulerHost:   config.SchedulerHost,
		SchedulerPath:   config.SchedulerPath,
		SchedulerPort:   config.SchedulerPort,
		SchedulerScheme: config.SchedulerScheme,
	}

	// marshal request data of websocket packet into json string
	reqDataString, err := json.Marshal(reqData)
	if err != nil {
		log.WithFields(log.Fields{
			"ClientID": streamer.ClientID,
			"Function": "SendSchedulerInfo",
		}).Warn("Failed to marshal request data into json string, abandon")
		return err
	}

	// send scheduler information
	streamer.Send(model.WSPacket{
		PacketType: "register_provider_metadata",
		Data:       string(reqDataString),
	}, nil)

	fmt.Printf("Sent scheduler info:\n%s\n", string(reqDataString))

	return nil
}

/*
	func: StateSelectedStorage
	description: state selected storage servers information to the provider streamer
*/
func (c *StreamerService) StateSelectedStorage(ctx context.Context, streamer *model.Streamer, instanceModel *model.InstanceModel) error {
	// define request format
	reqData := struct {
		StreamInstanceID string           `json:"stream_instance_id"`
		TargetFilestore  model.Filestore  `json:"selected_filestore"`
		TargetDepository model.Depository `json:"selected_depository"`
	}{
		StreamInstanceID: instanceModel.Instanceid,
		TargetFilestore:  instanceModel.TargetFilestore,
		TargetDepository: instanceModel.TargetDepository,
	}

	// marshal request data of websocket packet into json string
	reqDataString, err := json.Marshal(reqData)
	if err != nil {
		log.WithFields(log.Fields{
			"ClientID": streamer.ClientID,
			"Function": "StateSelectedStorage",
			"error":    err.Error(),
		}).Warn("Failed to marshal request data into json string, abandon")
		return err
	}

	// send selected storage servers information
	streamer.Send(model.WSPacket{
		PacketType: "state_selected_storage",
		Data:       string(reqDataString),
	}, nil)

	fmt.Printf("Sent selected storage:\n%s\n", string(reqDataString))

	return nil
}

/*
	func: StateNewInstance
	description: state the new wine container (instance) information to the provider streamer
*/
func (c *StreamerService) StateNewInstance(ctx context.Context, streamer *model.Streamer, instanceModel *model.InstanceModel) error {
	// 1. fill parameters in model.InstanceConnection

	var id_str string
	if instanceModel.VMID < 10 {
		id_str = "0" + strconv.Itoa(instanceModel.VMID)
	} else {
		id_str = strconv.Itoa(instanceModel.VMID % 100)
	}
	instanceModel.VideoRTCPort = "1" + id_str + "05"
	instanceModel.AudioRTCPort = "1" + id_str + "01"
	instanceModel.InputPort = "1" + id_str + "09"
	instanceModel.InstanceIP = "0.0.0.0"

	// 2. state VideoRTCPort & AudioRTCPort & InputPort

	// define request format
	reqData := instanceModel

	// marshal request data of websocket packet into json string
	reqDataString, err := json.Marshal(*reqData)
	if err != nil {
		log.WithFields(log.Fields{
			"ClientID": streamer.ClientID,
			"Function": "StateNewInstance",
			"error":    err.Error(),
		}).Warn("Failed to marshal request data into json string, abandon")
		return err
	}
	// send the new instance information
	streamer.Send(model.WSPacket{
		PacketType: "state_run_instance",
		Data:       string(reqDataString),
	}, nil)

	fmt.Printf("Sent instance info:\n%s\n", string(reqDataString))

	return nil
}

/*
	func: SendErrorMsg
	description: send error messages to the provider streamer
*/
func (c *StreamerService) SendErrorMsg(ctx context.Context, streamer *model.Streamer, errorType string, instanceID string) error {
	switch errorType {
	case model.ERROR_TYPE_STORAGE:
		reqData := struct {
			ErrorMessage     string `json:"error"`
			StreamInstanceID string `json:"stream_instance_id"`
		}{
			ErrorMessage:     "Failed to select storage.",
			StreamInstanceID: instanceID,
		}

		reqDataString, err := json.Marshal(reqData)
		if err != nil {
			log.WithFields(log.Fields{
				"ClientID": streamer.ClientID,
				"Function": "SendErrorMsg",
			}).Warn("Failed to marshal request data into json string, abandon")
			return err
		}

		streamer.Send(model.WSPacket{
			PacketType: "state_failed_select_storage",
			Data:       string(reqDataString),
		}, nil)
	case model.ERROR_TYPE_INSTANCE:
		reqData := struct {
			ErrorMessage     string `json:"error"`
			StreamInstanceID string `json:"stream_instance_id"`
		}{
			ErrorMessage:     "Failed to run instance.",
			StreamInstanceID: instanceID,
		}

		reqDataString, err := json.Marshal(reqData)
		if err != nil {
			log.WithFields(log.Fields{
				"ClientID": streamer.ClientID,
				"Function": "SendErrorMsg",
			}).Warn("Failed to marshal request data into json string, abandon")
			return err
		}

		streamer.Send(model.WSPacket{
			PacketType: "state_failed_run_instance",
			Data:       string(reqDataString),
		}, nil)
	case model.ERROR_TYPE_REMOVE:
		reqData := struct {
			ErrorMessage     string `json:"error"`
			StreamInstanceID string `json:"stream_instance_id"`
		}{
			ErrorMessage:     "Failed to remove instance.",
			StreamInstanceID: instanceID,
		}

		reqDataString, err := json.Marshal(reqData)
		if err != nil {
			log.WithFields(log.Fields{
				"ClientID": streamer.ClientID,
				"Function": "SendErrorMsg",
			}).Warn("Failed to marshal request data into json string, abandon")
			return err
		}

		streamer.Send(model.WSPacket{
			PacketType: "state_failed_remove_instance",
			Data:       string(reqDataString),
		}, nil)
	default:
		// leave empty
	}

	return nil
}
