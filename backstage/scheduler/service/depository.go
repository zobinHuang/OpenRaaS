package service

import (
	"context"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
@struct: DepositoryService
@description: service layer
*/
type DepositoryService struct {
	DepositoryDAL model.DepositoryDAL
}

/*
@struct: DepositoryServiceConfig
@description: used for config instance of struct DepositoryService
*/
type DepositoryServiceConfig struct {
	DepositoryDAL model.DepositoryDAL
}

func NewDepositoryService(f *DepositoryServiceConfig) model.DepositoryService {
	return &DepositoryService{
		DepositoryDAL: f.DepositoryDAL,
	}
}

func (f *DepositoryService) CreateDepositoryInRDS(ctx context.Context, info *model.DepositoryCore) error {
	return f.DepositoryDAL.CreateDepositoryInRDS(ctx, info)
}
