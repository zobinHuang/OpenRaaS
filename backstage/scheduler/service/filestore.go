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
	log.Infof("%s, allow new filestore enter, id: %s", utils.GetCurrentTime(), fileStore.ID)
	performance := "normal"
	if fileStore.IsContainFastNetspeed {
		performance = "powerful"
	}
	log.Infof("%s, New filestore id: %s, ip: %s, mem: %f GB, type: %s",
		utils.GetCurrentTime(), fileStore.ID, fileStore.IP, fileStore.Mem, performance)
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
	for _, f := range fileStores {
		if f.IsContainFastNetspeed {
			powerNum += 1
		} else {
			normalNum += 1
		}
		totalMem += f.Mem
	}
	log.Infof("%s, filestores info, Total: %d nodes, %f GB MEM, %d powerful node, %d normal node",
		utils.GetCurrentTime(), len(fileStores), totalMem, powerNum, normalNum)
	for _, f := range fileStores {
		performance := "normal"
		if f.IsContainFastNetspeed {
			performance = "powerful"
		}
		log.Infof("%s, filestore id: %s, ip: %s, mem: %f GB, type: %s",
			utils.GetCurrentTime(), f.ID, f.IP, f.Mem, performance)
	}
}
