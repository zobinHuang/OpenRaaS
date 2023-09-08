package servicecore

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
@struct: ScheduleServiceCore
@description: service core layer
*/
type ScheduleServiceCore struct {
	ConsumerDAL     model.ConsumerDAL
	ProviderDAL     model.ProviderDAL
	DepositoryDAL   model.DepositoryDAL
	FileStoreDAL    model.FileStoreDAL
	InstanceRoomDAL model.InstanceRoomDAL
	ApplicationDAL  model.ApplicationDAL
}

/*
@struct: ScheduleServiceCoreConfig
@description: used for config instance of struct ScheduleServiceCore
*/
type ScheduleServiceCoreConfig struct {
	ConsumerDAL     model.ConsumerDAL
	ProviderDAL     model.ProviderDAL
	DepositoryDAL   model.DepositoryDAL
	FileStoreDAL    model.FileStoreDAL
	InstanceRoomDAL model.InstanceRoomDAL
	ApplicationDAL  model.ApplicationDAL
}

/*
@func: NewScheduleServiceCore
@description:

	create, config and return an instance of struct ScheduleServiceCore
*/
func NewScheduleServiceCore(c *ScheduleServiceCoreConfig) model.ScheduleServiceCore {
	return &ScheduleServiceCore{
		ConsumerDAL:     c.ConsumerDAL,
		ProviderDAL:     c.ProviderDAL,
		DepositoryDAL:   c.DepositoryDAL,
		FileStoreDAL:    c.FileStoreDAL,
		InstanceRoomDAL: c.InstanceRoomDAL,
		ApplicationDAL:  c.ApplicationDAL,
	}
}

/*
@func: ScheduleStream
@description:

	core logic of scheduling stream instance is here
*/
func (sc *ScheduleServiceCore) ScheduleStream(ctx context.Context, streamInstance *model.StreamInstance) (*model.Provider, []model.DepositoryCore, []model.FileStoreCore, error) {
	providers := sc.ProviderDAL.GetProvider()
	appInfo, err := sc.ApplicationDAL.GetStreamApplicationByID(ctx, streamInstance.ApplicationID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scheduler GetStreamApplicationByID err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}

	candidatesGPU := make([]*model.Provider, 0, 0)
	if appInfo.IsProviderReqGPU {
		for _, p := range providers {
			if p.IsContainGPU {
				candidatesGPU = append(candidatesGPU, p)
			}
		}
	} else {
		candidatesGPU = providers
	}

	if len(candidatesGPU) <= 0 {
		return nil, nil, nil, fmt.Errorf("no provider can schedule")
	}

	if appInfo.FileStoreList == "" {
		return nil, nil, nil, fmt.Errorf("scheduler FileStoreList is none streamInstance: %+v", streamInstance)
	}
	var fileStoreStrList []string
	if err := json.Unmarshal([]byte(appInfo.FileStoreList), &fileStoreStrList); err != nil {
		return nil, nil, nil, fmt.Errorf("scheduler unmarshal FileStoreList fail, err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}

	// todo: get depositoryList from image id
	depositoryList, err := sc.DepositoryDAL.GetDepositoryInRDS(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scheduler GetDepositoryInRDS err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}

	filestoreList, err := sc.FileStoreDAL.GetFileStoreInRDSBetweenID(ctx, fileStoreStrList)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scheduler GetFileStoreInRDS err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}

	log.Infof("select info, provider: %+v, depositoryList: %+v, filestoreList: %+v", candidatesGPU[0], depositoryList, filestoreList)

	// Use Fisher-Yates algorithm to shuffle slices
	for i := len(filestoreList) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		filestoreList[i], filestoreList[j] = filestoreList[j], filestoreList[i]
	}

	for i := len(depositoryList) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		depositoryList[i], depositoryList[j] = depositoryList[j], depositoryList[i]
	}

	return candidatesGPU[0], depositoryList, filestoreList, nil
}

/*
@func: CreateStreamInstanceRoom
@description:

	create a room for the instance of stream instance
*/
func (sc *ScheduleServiceCore) CreateStreamInstanceRoom(ctx context.Context, provider *model.Provider,
	consumer *model.Consumer, streamInstance *model.StreamInstance) (*model.StreamInstanceRoom, error) {
	// initialize streamInstanceRoom instance
	streamInstanceRoom := &model.StreamInstanceRoom{
		StreamInstance: streamInstance,
		Provider:       provider,
	}

	// create consumer list, and insert our current consumer
	streamInstanceRoom.ConsumerList = make(map[string]*model.Consumer)
	streamInstanceRoom.ConsumerList[consumer.ClientID] = consumer

	// insert in dal layer
	sc.InstanceRoomDAL.CreateStreamInstanceRoom(ctx, streamInstanceRoom)

	return streamInstanceRoom, nil
}

// Clear delete all
func (sc *ScheduleServiceCore) Clear() {
	sc.ConsumerDAL.Clear()
	sc.ProviderDAL.Clear()
	sc.DepositoryDAL.Clear()
	sc.FileStoreDAL.Clear()
	sc.InstanceRoomDAL.Clear()
	sc.ApplicationDAL.Clear()
}
