package communicator

import (
	"context"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/zobinHuang/BrosCloud/provider/streamer/model"
)

/*
	@struct: WebsocketCommunicator
	@description: communicator to daemon and scheduler
*/
type WebsocketCommunicator struct {
	SchedulerWSConnection *model.Websocket
	DaemonWSConnection    *model.Websocket
	SchedulerDAL          model.SchedulerDAL
	DaemonDAL             model.DaemonDAL
}

/*
	@struct: WebsocketCommunicatorConfig
	@description: used for config instance of struct WebsocketCommunicator
*/
type WebsocketCommunicatorConfig struct {
	SchedulerDAL model.SchedulerDAL
	DaemonDAL    model.DaemonDAL
}

/*
	@func: NewWebsocketCommunicator
	@description:
		create, config and return an instance of struct WebsocketCommunicator
*/
func NewWebsocketCommunicator(c *WebsocketCommunicatorConfig) model.WebsocketCommunicator {
	wsCommunicator := &WebsocketCommunicator{
		SchedulerDAL: c.SchedulerDAL,
		DaemonDAL:    c.DaemonDAL,
	}

	wsCommunicator.InitSchdulerConnectionForDebug()

	return wsCommunicator
}

/*
	@func: InitSchdulerConnectionForDebug
	@description:
		initialize websocket communication just for debugging
*/
func (s *WebsocketCommunicator) InitSchdulerConnectionForDebug() {
	wsScheme := os.Getenv("SCHEDULER_WS_SCHEME")
	wsHostname := os.Getenv("SCHEDULER_WS_HOSTNAME")
	wsPort := os.Getenv("SCHEDULER_WS_PORT")
	wsPath := os.Getenv("SCHEDULER_WS_PATH")

	ctx := context.Background()

	// connect to scheduler
	err := s.ConnectToScheduler(ctx,
		wsScheme,
		wsHostname,
		wsPort,
		wsPath,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// start keep alive routine to scheduler
	s.KeepSchedulerConnAlive(ctx)

	// initilize scheduler recv callbacks route
	s.InitSchedulerRecvRoute(ctx)

	// generate request websocket to scheduler to register provider metadata
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
		}).Fatal("Failed to marshal provider metadata, abandoned")
	}

	// send metadata
	s.SchedulerWSConnection.Send(model.WSPacket{
		PacketType: "init_provider_metadata",
		Data:       string(reqToSchedulerString),
	}, nil)
}
