package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/utils"
)

/*
@struct: FileStoreService
@description: service layer
*/
type FileStoreService struct {
	FileStoreDAL model.FileStoreDAL
}

/*
@struct: FileStoreServiceConfig
@description: used for config instance of struct FileStoreService
*/
type FileStoreServiceConfig struct {
	FileStoreDAL model.FileStoreDAL
}

func NewFileStoreService(f *FileStoreServiceConfig) model.FileStoreService {
	return &FileStoreService{
		FileStoreDAL: f.FileStoreDAL,
	}
}

func (s *FileStoreService) CreateFileStoreInRDS(ctx context.Context, info *model.FileStoreCoreWithInst) error {
	return s.FileStoreDAL.CreateFileStoreInRDS(ctx, info)
}

func (s *FileStoreService) UpdateFileStoreInRDS(ctx context.Context, info *model.FileStoreCoreWithInst) error {
	return s.FileStoreDAL.UpdateFileStoreInRDSByID(ctx, info)
}

func (s *FileStoreService) ShowEnterInfo(ctx context.Context, fileStore *model.FileStoreCoreWithInst) {
	log.Infof("%s, 捕获到新节点上线, ID: %s", utils.GetCurrentTime(), fileStore.ID)
	log.Infof("认知到新的内容存储节点，详细信息：\n", fileStore.DetailedInfo())
}

func (s *FileStoreService) ShowAllInfo(ctx context.Context) {
	fileStores, err := s.FileStoreDAL.GetFileStoreInRDS(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("FileStoreService ShowAllInfo GetFileStoreInRDS error")
	}
	totalMem := 0.0
	powerNum := 0
	normalNum := 0
	abnormalNum := 0
	for _, f := range fileStores {
		if f.IsContainFastNetspeed {
			powerNum += 1
		} else {
			normalNum += 1
		}
		totalMem += f.Mem
		if f.GetAbnormalHistoryTimes() != 0 {
			abnormalNum += 1
		}
	}
	log.Infof("整合前，内容存储节点信息：%+v", fileStores)
	log.Infof("整合后，内容存储节点信息：")
	log.Infof("%s, 节点数量：%d, 总存储资源：%f GB MEM, 高性能节点数量：%d, 低性能节点数量：%d，异常节点数量： %d",
		utils.GetCurrentTime(), len(fileStores), totalMem, powerNum, normalNum, abnormalNum)
	log.Infof("ID 后有 * 表示服务异常节点")
	s.FileStoreDAL.ShowInfoFromRDS(fileStores)
}
