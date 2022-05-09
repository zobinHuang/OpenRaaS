package model

import (
	"context"

	"github.com/gorilla/websocket"
)

// --------- Service Layer Interface ---------
/*
	interface: RDbService
	description: interface of service layer for test
*/
type RDbService interface {
	GetRDbModel(ctx context.Context, rdm *RDbModel) error
}

/*
	interface: TokenService
	description: interface of token service for
			user authorization
*/
type TokenService interface {
	ValidateIDToken(tokenString string) (*User, error)
}

/*
	interface: ConsumerService
	description: interface of consumer service
*/
type ConsumerService interface {
	CreateConsumer(ctx context.Context, ws *websocket.Conn) (*Consumer, error)
	InitRecvRoute(ctx context.Context, consumer *Consumer)
}

/*
	interface: ApplicationService
	description: interface of application service
*/
type ApplicationService interface {
	GetStreamApplicationsCount(ctx context.Context) (int64, error)
	GetStreamApplicationDetails(ctx context.Context, applicationID string) (*StreamApplication, error)
	GetStreamApplications(ctx context.Context, pageNumber int, pageSize int, orderBy string) ([]*StreamApplication, error)
}

/*
	interface: ProviderService
	description: interface of provider service
*/
type ProviderService interface {
	CreateProvider(ctx context.Context, ws *websocket.Conn) (*Provider, error)
	InitRecvRoute(ctx context.Context, provider *Provider)
}

// --------- Service Core Layer Interface ---------
type ScheduleServiceCore interface {
	CreateStreamInstanceRoom(ctx context.Context, provider *Provider, consumer *Consumer, streamInstance *StreamInstance) (*StreamInstanceRoom, error)
	ScheduleStream(ctx context.Context, streamInstance *StreamInstance) (*Provider, []DepositaryCore, []FilestoreCore, error)
}

// --------- DAL Layer Interface ---------
/*
	interface: RDbDAL
	description: interface of data access layer for test
*/
type RDbDAL interface {
	GetRDbModel(ctx context.Context, rdm *RDbModel) error
}

/*
	interface: ConsumerDAL
	description: interface of data access layer for consumer
*/
type ConsumerDAL interface {
	CreateConsumer(ctx context.Context, consumer *Consumer)
	DeleteConsumer(ctx context.Context, consumerID string)
}

/*
	interface: ProviderDAL
	description: interface of data access layer for provider
*/
type ProviderDAL interface {
	CreateProvider(ctx context.Context, provider *Provider)
	DeleteProvider(ctx context.Context, providerID string)
}

/*
	interface: DepositaryDAL
	description: interface of data access layer for depositary
*/
type DepositaryDAL interface {
	CreateDepositary(ctx context.Context, depositary *Depositary)
	DeleteDepositary(ctx context.Context, depositaryID string)
}

/*
	interface: FilestoreDAL
	description: interface of data access layer for filestore
*/
type FilestoreDAL interface {
	CreateFilestore(ctx context.Context, filestore *Filestore)
	DeleteFilestore(ctx context.Context, filestoreID string)
}

/*
	interface: InstanceRoomDAL
	description: interface of data access layer for instance room
*/
type InstanceRoomDAL interface {
	CreateStreamInstanceRoom(ctx context.Context, streamInstanceRoom *StreamInstanceRoom)
	DeleteStreamInstanceRoom(ctx context.Context, instanceID string)
}

/*
	interface: ApplicationDAL
	description: interface of data access layer for application
*/
type ApplicationDAL interface {
	GetStreamApplicationsCount(ctx context.Context) (int64, error)
	GetStreamApplicationByID(ctx context.Context, applicationID string) (*StreamApplication, error)
	GetStreamApplicationsOrderedByUpdateTime(ctx context.Context, listLength int, listID int) ([]*StreamApplication, error)
	GetStreamApplicationsOrderedByName(ctx context.Context, listLength int, listID int) ([]*StreamApplication, error)
	GetStreamApplicationsOrderedByUsageCount(ctx context.Context, listLength int, listID int) ([]*StreamApplication, error)
}
