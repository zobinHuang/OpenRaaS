package communicator

import (
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

	return wsCommunicator
}
