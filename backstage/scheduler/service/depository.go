package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/utils"
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

func (s *DepositoryService) CreateDepositoryInRDS(ctx context.Context, info *model.DepositoryCoreWithInst) error {
	return s.DepositoryDAL.CreateDepositoryInRDS(ctx, info)
}

func (s *DepositoryService) UpdateFileStoreInRDS(ctx context.Context, info *model.DepositoryCoreWithInst) error {
	return s.DepositoryDAL.UpdateDepositoryInRDSByID(ctx, info)
}

func (s *DepositoryService) ShowEnterInfo(ctx context.Context, depository *model.DepositoryCoreWithInst) {
	log.Infof("%s, 捕获到新节点上线, ID: %s", utils.GetCurrentTime(), depository.ID)
	log.Infof("认知到新的镜像仓库节点，详细信息：%s\n", depository.DetailedInfo())
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
	abnormalNum := 0
	for _, d := range depositories {
		if d.IsContainFastNetspeed {
			powerNum += 1
		} else {
			normalNum += 1
		}
		totalMem += d.Mem
		if d.GetAbnormalHistoryTimes() != 0 {
			abnormalNum += 1
		}
	}
	log.Infof("整合前，镜像仓库节点信息：%+v", depositories)
	log.Infof("整合后，镜像仓库节点信息：")
	log.Infof("%s, 节点数量：%d, 总存储资源：%f GB MEM, 高性能节点数量：%d, 低性能节点数量：%d，异常节点数量： %d",
		utils.GetCurrentTime(), len(depositories), totalMem, powerNum, normalNum, abnormalNum)
	log.Infof("ID 后有 * 表示服务异常节点")
	s.DepositoryDAL.ShowInfoFromRDS(depositories)
}
