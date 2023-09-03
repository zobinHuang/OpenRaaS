package service

import (
	"context"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
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

func (f *FileStoreService) CreateFileStoreInRDS(ctx context.Context, info *model.FileStoreCore) error {
	return f.FileStoreDAL.CreateFileStoreInRDS(ctx, info)
}
