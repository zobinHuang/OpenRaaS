package service

import "github.com/zobinHuang/BrosCloud/provider/streamer/model"

/*
	@struct: InstanceService
	@description: service layer
*/
type InstanceService struct {
	InstanceDAL model.InstanceDAL
}

/*
	@struct: InstanceServiceConfig
	@description: used for config instance of struct InstanceService
*/
type InstanceServiceConfig struct {
	InstanceDAL model.InstanceDAL
}

/*
	@func: NewInstanceService
	@description:
		create, config and return an instance of struct InstanceService
*/
func NewInstanceService(c *InstanceServiceConfig) model.InstanceService {
	return &InstanceService{
		InstanceDAL: c.InstanceDAL,
	}
}
