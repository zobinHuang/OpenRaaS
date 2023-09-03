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
	Clear()
}

/*
interface: ApplicationService
description: interface of application service
*/
type ApplicationService interface {
	GetStreamApplicationsCount(ctx context.Context) (int64, error)
	GetStreamApplicationDetails(ctx context.Context, applicationID string) (*StreamApplication, error)
	GetStreamApplications(ctx context.Context, pageNumber int, pageSize int, orderBy string) ([]*StreamApplication, error)
	CreateStreamApplication(ctx context.Context, info *StreamApplication) error
}

/*
interface: ProviderService
description: interface of provider service
*/
type ProviderService interface {
	CreateProvider(ctx context.Context, ws *websocket.Conn, uuid string) (*Provider, error)
	InitRecvRoute(ctx context.Context, provider *Provider)
	CreateProviderInRDS(ctx context.Context, provider *ProviderCore) error
}

/*
interface: FileStoreService
description: interface of FileStore service
*/
type FileStoreService interface {
	CreateFileStoreInRDS(ctx context.Context, info *FileStoreCore) error
}

/*
interface: DepositoryService
description: interface of Depository service
*/
type DepositoryService interface {
	CreateDepositoryInRDS(ctx context.Context, info *DepositoryCore) error
}

// --------- Service Core Layer Interface ---------
type ScheduleServiceCore interface {
	CreateStreamInstanceRoom(ctx context.Context, provider *Provider, consumer *Consumer, streamInstance *StreamInstance) (*StreamInstanceRoom, error)
	ScheduleStream(ctx context.Context, streamInstance *StreamInstance) (*Provider, []DepositoryCore, []FileStoreCore, error)
	Clear()
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
	Clear()
}

/*
interface: ProviderDAL
description: interface of data access layer for provider
*/
type ProviderDAL interface {
	CreateProvider(ctx context.Context, provider *Provider)
	DeleteProvider(ctx context.Context, providerID string)
	GetProvider() []*Provider
	GetProviderInRDS(ctx context.Context) ([]ProviderCore, error)
	GetProviderInRDSByID(ctx context.Context, id string) (*ProviderCore, error)
	DeleteProviderInRDSByID(ctx context.Context, id string) error
	UpdateProviderInRDSByID(ctx context.Context, provider *ProviderCore) error
	CreateProviderInRDS(ctx context.Context, provider *ProviderCore) error
	Clear()
}

/*
interface: DepositoryDAL
description: interface of data access layer for depositary
*/
type DepositoryDAL interface {
	CreateDepository(ctx context.Context, depositary *Depository)
	DeleteDepository(ctx context.Context, depositaryID string)
	GetDepositoryInRDS(ctx context.Context) ([]DepositoryCore, error)
	GetDepositoryBetweenIDInRDS(ctx context.Context, ids []string) ([]DepositoryCore, error)
	GetDepositoryInRDSByID(ctx context.Context, id string) (*DepositoryCore, error)
	DeleteDepositoryInRDSByID(ctx context.Context, id string) error
	UpdateDepositoryInRDSByID(ctx context.Context, info *DepositoryCore) error
	CreateDepositoryInRDS(ctx context.Context, info *DepositoryCore) error
	Clear()
}

/*
interface: FileStoreDAL
description: interface of data access layer for filestore
*/
type FileStoreDAL interface {
	CreateFileStore(ctx context.Context, fileStore *FileStore)
	DeleteFileStore(ctx context.Context, id string)
	GetFileStoreInRDS(ctx context.Context) ([]FileStoreCore, error)
	GetFileStoreInRDSByID(ctx context.Context, id string) (*FileStoreCore, error)
	GetFileStoreInRDSBetweenID(ctx context.Context, ids []string) ([]FileStoreCore, error)
	DeleteFileStoreInRDSByID(ctx context.Context, id string) error
	UpdateFileStoreInRDSByID(ctx context.Context, info *FileStoreCore) error
	CreateFileStoreInRDS(ctx context.Context, info *FileStoreCore) error
	Clear()
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
	Clear()
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
	Clear()
}
