package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/utils"
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

func (s *DepositoryService) CreateDepositoryInRDS(ctx context.Context, info *model.DepositoryCore) error {
	return s.DepositoryDAL.CreateDepositoryInRDS(ctx, info)
}

func (s *DepositoryService) ShowEnterInfo(ctx context.Context, depository *model.DepositoryCore) {
	log.Infof("%s, allow new depository enter, id: %s", utils.GetCurrentTime(), depository.ID)
	performance := "normal"
	if depository.IsContainFastNetspeed {
		performance = "powerful"
	}
	log.Infof("%s, New depository id: %s, ip: %s, mem: %f GB, type: %s",
		utils.GetCurrentTime(), depository.ID, depository.IP, depository.Mem, performance)
}

func (s *DepositoryService) ShowAllInfo(ctx context.Context) {
	depositories, err := s.DepositoryDAL.GetDepositoryInRDS(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("DepositoryService ShowAllInfo GetDepositoryInRDS error")
	}
	totalMem := 0.0
	powerNum := 0
	normalNum := 0
	for _, d := range depositories {
		if d.IsContainFastNetspeed {
			powerNum += 1
		} else {
			normalNum += 1
		}
		totalMem += d.Mem
	}
	log.Infof("%s, Depositories Info, Total: %d nodes, %f GB MEM, %d powerful node, %d normal node",
		utils.GetCurrentTime(), len(depositories), totalMem, powerNum, normalNum)
	for _, d := range depositories {
		performance := "normal"
		if d.IsContainFastNetspeed {
			performance = "powerful"
		}
		log.Infof("%s, depository id: %s, ip: %s, mem: %f GB, type: %s",
			utils.GetCurrentTime(), d.ID, d.IP, d.Mem, performance)
	}
}
