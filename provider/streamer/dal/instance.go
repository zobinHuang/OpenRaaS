package dal

import (
	"context"
	"fmt"

	"github.com/zobinHuang/BrosCloud/provider/streamer/model"
)

/*
	@struct: InstanceDAL
	@description: DAL layer
*/
type InstanceDAL struct {
	StreamInstanceList map[string]*model.StreamInstanceDaemonModel
}

/*
	@struct: InstanceDALConfig
	@description: used for config instance of struct InstanceDAL
*/
type InstanceDALConfig struct {
}

/*
	@function: NewInstanceDAL
	@description:
		create, config and return an instance of struct InstanceDAL
*/
func NewInstanceDAL(c *InstanceDALConfig) model.InstanceDAL {
	idal := &InstanceDAL{}

	idal.StreamInstanceList = make(map[string]*model.StreamInstanceDaemonModel)

	return idal
}

/*
	@function: AddNewStreamInstance
	@description:
		append new stream instance
*/
func (idal *InstanceDAL) AddNewStreamInstance(ctx context.Context, streamInstanceDaemonModel *model.StreamInstanceDaemonModel) {
	idal.StreamInstanceList[streamInstanceDaemonModel.Instanceid] = streamInstanceDaemonModel
}

/*
	@function: GetStreamInstanceByID
	@description:
		obtain stream instance based on given instance index
*/
func (idal *InstanceDAL) GetStreamInstanceByID(ctx context.Context, streamInstanceID string) (*model.StreamInstanceDaemonModel, error) {
	streamInstance, ok := idal.StreamInstanceList[streamInstanceID]
	if !ok {
		return nil, fmt.Errorf("Can't find stream instance with id %s\n", streamInstanceID)
	}

	return streamInstance, nil
}
