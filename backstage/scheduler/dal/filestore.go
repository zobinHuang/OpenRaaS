package dal

import (
	"context"
	"gorm.io/gorm"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
@struct: FileStoreDAL
@description: DAL layer
*/
type FileStoreDAL struct {
	DB *gorm.DB
}

/*
@struct: FileStoreDALConfig
@description: used for config instance of struct FileStoreDAL
*/
type FileStoreDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewFileStoreDAL
@description:

	create, config and return an instance of struct FileStoreDAL
*/
func NewFileStoreDAL(c *FileStoreDALConfig) model.FileStoreDAL {
	return &FileStoreDAL{
		DB: c.DB,
	}
}

/*
@func: CreateFileStore
@description:

	insert a new filestore to FileStore list
*/
func (d *FileStoreDAL) CreateFileStore(ctx context.Context, info *model.FileStore) error {
	// todo
	return nil
}

/*
@func: DeleteFileStore
@description:

	delete the specified filestore from FileStore list
*/
func (d *FileStoreDAL) DeleteFileStoreByID(ctx context.Context, id string) error {
	return nil
}

func (d *FileStoreDAL) GetFileStore(ctx context.Context) ([]model.FileStore, error) {
	return nil, nil
}
func (d *FileStoreDAL) GetFileStoreByID(ctx context.Context, id string) (*model.FileStore, error) {
	return nil, nil
}
func (d *FileStoreDAL) UpdateFileStoreByID(ctx context.Context, info *model.FileStore) error {
	return nil
}
