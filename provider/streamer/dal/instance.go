package dal

import "github.com/zobinHuang/BrosCloud/provider/streamer/model"

type InstanceDAL struct {
	StreamInstanceList map[string]*model.StreamInstance
}

type InstanceDALConfig struct {
}

func NewInstanceDAL(c *InstanceDALConfig) model.InstanceDAL {
	idal := &InstanceDAL{}

	idal.StreamInstanceList = make(map[string]*model.StreamInstance)

	return idal
}
