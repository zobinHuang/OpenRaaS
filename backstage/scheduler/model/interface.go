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
	GetScheduleServiceCore() *ScheduleServiceCore
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
	AddFileStoreIDToAPPInRDS(ctx context.Context, info *StreamApplication, id string) error
	ShowEnterInfo(ctx context.Context, app *StreamApplication, nodeId string)
	ShowAllInfo(ctx context.Context)
}

/*
interface: ProviderService
description: interface of provider service
*/
type ProviderService interface {
	CreateProvider(ctx context.Context, ws *websocket.Conn, uuid string) (*Provider, error)
	InitRecvRoute(ctx context.Context, provider *Provider)
	CreateProviderInRDS(ctx context.Context, provider *ProviderCoreWithInst) error
	UpdateProviderInRDS(ctx context.Context, provider *ProviderCoreWithInst) error
	ShowEnterInfo(ctx context.Context, provider *ProviderCoreWithInst)
	ShowAllInfo(ctx context.Context)
}

/*
interface: FileStoreService
description: interface of FileStore service
*/
type FileStoreService interface {
	CreateFileStoreInRDS(ctx context.Context, info *FileStoreCoreWithInst) error
	UpdateFileStoreInRDS(ctx context.Context, info *FileStoreCoreWithInst) error
	ShowEnterInfo(ctx context.Context, fileStore *FileStoreCoreWithInst)
	ShowAllInfo(ctx context.Context)
}

/*
interface: DepositoryService
description: interface of Depository service
*/
type DepositoryService interface {
	CreateDepositoryInRDS(ctx context.Context, info *DepositoryCoreWithInst) error
	UpdateFileStoreInRDS(ctx context.Context, info *DepositoryCoreWithInst) error
	ShowEnterInfo(ctx context.Context, depository *DepositoryCoreWithInst)
	ShowAllInfo(ctx context.Context)
}

// --------- Service Core Layer Interface ---------
type ScheduleServiceCore interface {
	CreateStreamInstanceRoom(ctx context.Context, provider *Provider, consumer *Consumer, streamInstance *StreamInstance) (*StreamInstanceRoom, error)
	ScheduleStream(ctx context.Context, consumer *Consumer, streamInstance *StreamInstance) (*Provider, []DepositoryCoreWithInst, []FileStoreCoreWithInst, error)
	GetStreamInstanceRoomByInstanceID(id string) (*StreamInstanceRoom, error)
	SetValueToBlockchain(key, value string) error
	GetValueFromBlockchain(key string) (string, error)
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
	GetProviderInRDS(ctx context.Context) ([]ProviderCoreWithInst, error)
	GetProviderInRDSByID(ctx context.Context, id string) (*ProviderCoreWithInst, error)
	DeleteProviderInRDSByID(ctx context.Context, id string) error
	UpdateProviderInRDSByID(ctx context.Context, provider *ProviderCoreWithInst) error
	CreateProviderInRDS(ctx context.Context, provider *ProviderCoreWithInst) error
	Clear()
}

/*
interface: DepositoryDAL
description: interface of data access layer for depository
*/
type DepositoryDAL interface {
	CreateDepository(ctx context.Context, depository *Depository)
	DeleteDepository(ctx context.Context, depositoryID string)
	GetDepositoryInRDS(ctx context.Context) ([]DepositoryCoreWithInst, error)
	GetDepositoryBetweenIDInRDS(ctx context.Context, ids []string) ([]DepositoryCoreWithInst, error)
	GetDepositoryInRDSByID(ctx context.Context, id string) (*DepositoryCoreWithInst, error)
	DeleteDepositoryInRDSByID(ctx context.Context, id string) error
	UpdateDepositoryInRDSByID(ctx context.Context, info *DepositoryCoreWithInst) error
	CreateDepositoryInRDS(ctx context.Context, info *DepositoryCoreWithInst) error
	Clear()
}

/*
interface: FileStoreDAL
description: interface of data access layer for filestore
*/
type FileStoreDAL interface {
	CreateFileStore(ctx context.Context, fileStore *FileStore)
	DeleteFileStore(ctx context.Context, id string)
	GetFileStoreInRDS(ctx context.Context) ([]FileStoreCoreWithInst, error)
	GetFileStoreInRDSByID(ctx context.Context, id string) (*FileStoreCoreWithInst, error)
	GetFileStoreInRDSBetweenID(ctx context.Context, ids []string) ([]FileStoreCoreWithInst, error)
	DeleteFileStoreInRDSByID(ctx context.Context, id string) error
	UpdateFileStoreInRDSByID(ctx context.Context, info *FileStoreCoreWithInst) error
	CreateFileStoreInRDS(ctx context.Context, info *FileStoreCoreWithInst) error
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
	GetInstanceRoomByInstanceID(ctx context.Context, instanceID string) (*StreamInstanceRoom, error)
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
