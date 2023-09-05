package service

import (
	"context"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
	struct: rdService
	description: service layer
*/
type rdbService struct {
	RDbDAL model.RDbDAL
}

/*
	struct: RDbServiceConfig
	description: used for config instance of struct rdbService
*/
type RDbServiceConfig struct {
	RDbDAL model.RDbDAL
}

/*
	func: NewRDbService
	description: create, config and return an instance of struct rdbService
*/
func NewRDbService(c *RDbServiceConfig) model.RDbService {
	return &rdbService{
		RDbDAL: c.RDbDAL,
	}
}

/*
	func: GetRDbModel
	description: service that return err message for test
*/
func (s *rdbService) GetRDbModel(ctx context.Context, rdbm *model.RDbModel) error {
	err := s.RDbDAL.GetRDbModel(ctx, rdbm)
	return err
}
