package service

import (
	"context"

	"github.com/zobinHuang/BrosCloud/provider/streamer/model"
)

/*
	@struct: SchedulerService
	@description: service layer
*/
type SchedulerService struct {
	SchedulerDAL          model.SchedulerDAL
	WebsocketCommunicator model.WebsocketCommunicator
}

/*
	@struct: SchedulerServiceConfig
	@description: used for config instance of struct SchedulerService
*/
type SchedulerServiceConfig struct {
	SchedulerDAL          model.SchedulerDAL
	WebsocketCommunicator model.WebsocketCommunicator
}

/*
	@func: NewSchedulerService
	@description:
		create, config and return an instance of struct SchedulerService
*/
func NewSchedulerService(c *SchedulerServiceConfig) model.SchedulerService {
	return &SchedulerService{
		SchedulerDAL:          c.SchedulerDAL,
		WebsocketCommunicator: c.WebsocketCommunicator,
	}
}

/*
	@func: ConnectToScheduler
	@description:
		connect to scheduler node
*/
func (s *SchedulerService) ConnectToScheduler(ctx context.Context, scheme string, hostname string, port string, path string) error {
	return s.WebsocketCommunicator.ConnectToScheduler(ctx, scheme, hostname, port, path)
}

/*
	@func: KeepSchedulerConnAlive
	@description:
		keep alive routine
*/
func (s *SchedulerService) KeepSchedulerConnAlive(ctx context.Context) {
	s.WebsocketCommunicator.KeepSchedulerConnAlive(ctx)
}

/*
	@func: InitSchedulerRecvRoute
	@description:
		initialize receiving callback
*/
func (s *SchedulerService) InitSchedulerRecvRoute(ctx context.Context) {
	s.WebsocketCommunicator.InitSchedulerRecvRoute(ctx)
}
