package communicator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/zobinHuang/OpenRaaS/provider/streamer/model"
)

/*
	@struct: WebsocketCommunicator
	@description: communicator to daemon and scheduler
*/
type WebsocketCommunicator struct {
	SchedulerWSConnection *model.Websocket
	DaemonWSConnection    *model.Websocket
	InstanceDAL           model.InstanceDAL
	SchedulerDAL          model.SchedulerDAL
	DaemonDAL             model.DaemonDAL
	WebRTCStreamDAL       model.WebRTCStreamDAL
}

/*
	@struct: WebsocketCommunicatorConfig
	@description: used for config instance of struct WebsocketCommunicator
*/
type WebsocketCommunicatorConfig struct {
	InstanceDAL     model.InstanceDAL
	SchedulerDAL    model.SchedulerDAL
	DaemonDAL       model.DaemonDAL
	WebRTCStreamDAL model.WebRTCStreamDAL
}

/*
	@func: NewWebsocketCommunicator
	@description:
		create, config and return an instance of struct WebsocketCommunicator
*/
func NewWebsocketCommunicator(c *WebsocketCommunicatorConfig) model.WebsocketCommunicator {
	wsCommunicator := &WebsocketCommunicator{
		SchedulerDAL:    c.SchedulerDAL,
		DaemonDAL:       c.DaemonDAL,
		InstanceDAL:     c.InstanceDAL,
		WebRTCStreamDAL: c.WebRTCStreamDAL,
	}

	err := wsCommunicator.InitDaemonConnection()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Failed to build connection to provider daemon")
	}

	// wsCommunicator.InitSchdulerConnectionForDebug()

	return wsCommunicator
}

func (s *WebsocketCommunicator) InitDaemonConnection() error {
	ctx := context.Background()

	// obtain websocket metadata to daemon
	wsScheme := os.Getenv("DAEMON_WS_SCHEME")
	wsHostname := os.Getenv("DAEMON_WS_HOSTNAME")
	wsPort := os.Getenv("DAEMON_WS_PORT")
	wsPath := os.Getenv("DAEMON_WS_PATH")

	completeHostname := fmt.Sprintf("%s:%s", wsHostname, wsPort)
	daemonURL := url.URL{
		Scheme:   wsScheme,
		Host:     completeHostname,
		Path:     wsPath,
		RawQuery: "type=provider",
	}

	conn, _, err := websocket.DefaultDialer.Dial(daemonURL.String(), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to build websocket connection to daemon")
		return err
	}

	// store websocket to daemon, and start to serve it
	s.NewDaemonConnection(ctx, conn)

	// register recv callbacks
	s.InitDaemonRecvRoute(ctx)

	// start keep alive
	s.KeepDaemonConnAlive(ctx)

	return nil
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
