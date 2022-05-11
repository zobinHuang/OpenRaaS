package model

import (
	"context"

	"github.com/gorilla/websocket"
)

/*
	interface: InstanceService
	description: interface of service layer for instance
*/
type InstanceService interface {
}

/*
	interface: SchedulerService
	description: interface of service layer for scheduler
*/
type SchedulerService interface {
	ConnectToScheduler(ctx context.Context, scheme string, hostname string, port string, path string) error
	KeepSchedulerConnAlive(ctx context.Context)
	InitSchedulerRecvRoute(ctx context.Context)
}

/*
	interface: DaemonService
	description: interface of service layer for daemon
*/
type DaemonService interface {
}

/*
	interface: WebsocketCommunicator
	description: interface of service layer for websocket
*/
type WebsocketCommunicator interface {
	// communicator to scheduler
	ConnectToScheduler(ctx context.Context, scheme string, hostname string, port string, path string) error
	KeepSchedulerConnAlive(ctx context.Context)
	InitSchedulerRecvRoute(ctx context.Context)

	// communicator to daemon
	NewDaemonConnection(ctx context.Context, conn *websocket.Conn)
	InitDaemonRecvRoute(ctx context.Context)
}

/*
	interface: InstanceDAL
	description: interface of data access layer for instance
*/
type InstanceDAL interface {
	AddNewStreamInstance(ctx context.Context, streamInstanceDaemonModel *StreamInstanceDaemonModel)
	GetStreamInstanceByID(ctx context.Context, streamInstanceID string) (*StreamInstanceDaemonModel, error)
}

/*
	interface: SchedulerDAL
	description: interface of data access layer for scheduler
*/
type SchedulerDAL interface {
	SetICEServers(iceServer string)
}

/*
	interface: DaemonDAL
	description: interface of data access layer for daemon
*/
type DaemonDAL interface {
}

/*
	interface: WebRTCStreamDAL
	description: interface of data access layer for webrtc streamer
*/
type WebRTCStreamDAL interface {
	NewWebRTCStreamer(ctx context.Context, streamInstance *StreamInstanceDaemonModel) (*WebRTCStreamer, error)
}
