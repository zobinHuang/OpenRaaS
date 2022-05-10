package dal

import (
	"github.com/zobinHuang/BrosCloud/provider/streamer/model"
)

type InstanceDAL struct {
	StreamInstanceList map[string]*model.StreamInstanceDaemonModel
}

type InstanceDALConfig struct {
}

func NewInstanceDAL(c *InstanceDALConfig) model.InstanceDAL {
	idal := &InstanceDAL{}

	idal.StreamInstanceList = make(map[string]*model.StreamInstanceDaemonModel)

	return idal
}

func (idal *InstanceDAL) AddNewStreamInstance(streamInstanceDaemonModel *model.StreamInstanceDaemonModel) {
	idal.StreamInstanceList[streamInstanceDaemonModel.Instanceid] = streamInstanceDaemonModel
}
