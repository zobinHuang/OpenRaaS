package model

import (
	"context"

	"github.com/gorilla/websocket"
)

/*
	interface: RDbService
	description: interface of service layer for test
*/
type RDbService interface {
	GetRDbModel(ctx context.Context, rdm *RDbModel) error
}

/*
	interface: RDbDAL
	description: interface of data access layer for test
*/
type RDbDAL interface {
	GetRDbModel(ctx context.Context, rdm *RDbModel) error
}

// --------- Service Layer Interface ---------
/*
	interface: InstanceService
	description: interface of instance, which includes some func for wine docker management
*/
type InstanceService interface {
	NewVMID(ctx context.Context, instanceModel *InstanceModel) error
	LaunchInstance(ctx context.Context, instanceModel *InstanceModel) chan struct{}
	DeleteInstance(ctx context.Context, vmid int) error
	DeleteInstanceByInstanceid(ctx context.Context, Instanceid string) error
	DeleteAllInstance(ctx context.Context) error
	MountFilestore(ctx context.Context, vmid int, filestore FilestoreCore) error
	FetchLayerFromDepository(ctx context.Context, vmid int, depository DepositoryCore, imageName string) error
}

/*
	interface: StreamerService
	description: interface of provider streamer, used to start or contact with provider streamer
*/
type StreamerService interface {
	RunStreamerContainer(ctx context.Context) error
	KillStreamerContainer(ctx context.Context) error
	CreateStreamer(ctx context.Context, ws *websocket.Conn) (*Streamer, error)
	StateScheduler(ctx context.Context, streamer *Streamer) error
	StateSelectedStorage(ctx context.Context, streamer *Streamer, instanceModel *InstanceModel) error
	StateNewInstance(ctx context.Context, streamer *Streamer, instanceModel *InstanceModel) error
	SendErrorMsg(ctx context.Context, streamer *Streamer, errorType string, instanceID string) error
}
