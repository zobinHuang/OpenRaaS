package service

import "github.com/zobinHuang/OpenRaaS/provider/streamer/model"

type DaemonService struct {
	DaemonDAL             model.DaemonDAL
	WebsocketCommunicator model.WebsocketCommunicator
}

type DaemonServiceConfig struct {
	DaemonDAL             model.DaemonDAL
	WebsocketCommunicator model.WebsocketCommunicator
}

func NewDaemonService(c *DaemonServiceConfig) model.DaemonService {
	return &DaemonService{
		DaemonDAL:             c.DaemonDAL,
		WebsocketCommunicator: c.WebsocketCommunicator,
	}
}
