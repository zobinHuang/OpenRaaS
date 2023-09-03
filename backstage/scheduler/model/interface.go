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

/*
interface: FileStoreService
description: interface of FileStore service
*/
type FileStoreService interface {
	CreateFileStore(ctx context.Context, ws *websocket.Conn) (*FileStore, error)
	InitRecvRoute(ctx context.Context, provider *FileStore)
}

/*
interface: DepositaryService
description: interface of Depositary service
*/
type DepositaryService interface {
	CreateDepositary(ctx context.Context, ws *websocket.Conn) (*FileStore, error)
	InitRecvRoute(ctx context.Context, provider *FileStore)
}

// --------- Service Core Layer Interface ---------
type ScheduleServiceCore interface {
	CreateStreamInstanceRoom(ctx context.Context, provider *Provider, consumer *Consumer, streamInstance *StreamInstance) (*StreamInstanceRoom, error)
	ScheduleStream(ctx context.Context, streamInstance *StreamInstance) (*Provider, []DepositaryCore, []FileStoreCore, error)
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
	GetConsumerByID(ctx context.Context, consumerID string) (*Consumer, error)
}

/*
interface: ProviderDAL
description: interface of data access layer for provider
*/
type ProviderDAL interface {
	GetProvider(ctx context.Context) ([]Provider, error)
	GetProviderByID(ctx context.Context, id string) (*Provider, error)
	DeleteProviderByID(ctx context.Context, id string) error
	UpdateProviderByID(ctx context.Context, info *Provider) error
	CreateProvider(ctx context.Context, info *Provider) error
}

/*
interface: DepositaryDAL
description: interface of data access layer for depositary
*/
type DepositaryDAL interface {
	GetDepositary(ctx context.Context) ([]Depositary, error)
	GetDepositaryByID(ctx context.Context, id string) (*Depositary, error)
	DeleteDepositaryByID(ctx context.Context, id string) error
	UpdateDepositaryByID(ctx context.Context, info *Depositary) error
	CreateDepositary(ctx context.Context, info *Depositary) error
}

/*
interface: FileStoreDAL
description: interface of data access layer for filestore
*/
type FileStoreDAL interface {
	GetFileStore(ctx context.Context) ([]FileStore, error)
	GetFileStoreByID(ctx context.Context, id string) (*FileStore, error)
	DeleteFileStoreByID(ctx context.Context, id string) error
	UpdateFileStoreByID(ctx context.Context, info *FileStore) error
	CreateFileStore(ctx context.Context, info *FileStore) error
}

/*
interface: InstanceRoomDAL
description: interface of data access layer for instance room
*/
type InstanceRoomDAL interface {
	CreateStreamInstanceRoom(ctx context.Context, streamInstanceRoom *StreamInstanceRoom)
	DeleteStreamInstanceRoom(ctx context.Context, instanceID string)
	GetConsumerMapByInstanceID(ctx context.Context, instanceID string) (map[string]*Consumer, error)
	GetProviderByInstanceID(ctx context.Context, instanceID string) (*Provider, error)
}

/*
interface: ApplicationDAL
description: interface of data access layer for application
*/
type ApplicationDAL interface {
	GetStreamApplicationsCount(ctx context.Context) (int64, error)
	GetStreamApplication(ctx context.Context) ([]StreamApplication, error)
	GetStreamApplicationByID(ctx context.Context, applicationID string) (*StreamApplication, error)
	GetStreamApplicationsOrderedByUpdateTime(ctx context.Context, listLength int, listID int) ([]*StreamApplication, error)
	GetStreamApplicationsOrderedByName(ctx context.Context, listLength int, listID int) ([]*StreamApplication, error)
	GetStreamApplicationsOrderedByUsageCount(ctx context.Context, listLength int, listID int) ([]*StreamApplication, error)
	DeleteStreamApplicationByID(ctx context.Context, id string) error
	UpdateStreamApplicationByID(ctx context.Context, info *StreamApplication) error
	CreateStreamApplication(ctx context.Context, info *StreamApplication) error
}
