package dal

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
@struct: FileStoreDAL
@description: DAL layer
*/
type FileStoreDAL struct {
	DB            *gorm.DB
	FileStoreList map[string]*model.FileStore
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
	ddal := &FileStoreDAL{}

	ddal.FileStoreList = make(map[string]*model.FileStore)
	ddal.DB = c.DB

	return ddal
}

/*
@func: CreateFileStore
@description:

	insert a new filestore to FileStore list
*/
func (d *FileStoreDAL) CreateFileStore(ctx context.Context, filestore *model.FileStore) {
	d.FileStoreList[filestore.ClientID] = filestore
}

/*
@func: DeleteFileStore
@description:

	delete the specified filestore from FileStore list
*/
func (d *FileStoreDAL) DeleteFileStore(ctx context.Context, filestoreID string) {
	delete(d.FileStoreList, filestoreID)
}

// CreateFileStoreInRDS create file store core info to rds
func (d *FileStoreDAL) CreateFileStoreInRDS(ctx context.Context, info *model.FileStoreCore) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_cores").Create(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"info":  info,
		}).Warn("Failed to create file store core info to rds")
		return err
	}

	return nil
}

// DeleteFileStoreInRDSByID delete file store core info by id in rds
func (d *FileStoreDAL) DeleteFileStoreInRDSByID(ctx context.Context, id string) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_cores").Where("id=?", id).Delete(&model.FileStoreCore{}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to delete file store core info by id in rds")
		return err
	}

	return nil
}

// GetFileStoreInRDS obtain all file store core info from rds
func (d *FileStoreDAL) GetFileStoreInRDS(ctx context.Context) ([]model.FileStoreCore, error) {
	var infos []model.FileStoreCore

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_cores").Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all file store core info from rds")
		return nil, err
	}

	return infos, nil
}

// GetFileStoreInRDSByID get file store core info by id from rds
func (d *FileStoreDAL) GetFileStoreInRDSByID(ctx context.Context, id string) (*model.FileStoreCore, error) {
	var info model.FileStoreCore
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_cores").Where("id = ?", id).First(&info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to get file store core info by id from rds")
		return nil, err
	}
	return &info, nil
}

// UpdateFileStoreInRDSByID update file store core info by id in rds
func (d *FileStoreDAL) UpdateFileStoreInRDSByID(ctx context.Context, info *model.FileStoreCore) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_cores").Where("id=?", info.ID).Updates(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    info.ID,
		}).Warn("Failed to update file store core info by id in rds")
		return err
	}
	return nil
}

// GetFileStoreInRDSBetweenID get file store core info between id from rds
func (d *FileStoreDAL) GetFileStoreInRDSBetweenID(ctx context.Context, ids []string) ([]model.FileStoreCore, error) {
	var infos []model.FileStoreCore

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_cores").Where("id IN (?)", ids).Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all file store core info from rds")
		return nil, err
	}

	return infos, nil
}

// Clear delete all
func (d *FileStoreDAL) Clear() {
	if err := d.DB.Exec("DELETE FROM public.file_store_cores").Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Fail to clear file_store_cores table")
	}
}
