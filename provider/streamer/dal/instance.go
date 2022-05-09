package dal

import "github.com/zobinHuang/BrosCloud/provider/streamer/model"

type InstanceDAL struct {
	InstanceList map[string]*model.Instance
}

type InstanceDALConfig struct {
}

func NewInstanceDAL(c *InstanceDALConfig) model.InstanceDAL {
	idal := &InstanceDAL{}

	idal.InstanceList = make(map[string]*model.Instance)

	return idal
}
