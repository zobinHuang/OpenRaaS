package dal

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
@struct: DepositoryDAL
@description: DAL layer
*/
type DepositoryDAL struct {
	DB             *gorm.DB
	DepositoryList map[string]*model.Depository
}

/*
@struct: DepositoryDALConfig
@description: used for config instance of struct DepositoryDAL
*/
type DepositoryDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewDepositoryDAL
@description:

	create, config and return an instance of struct DepositoryDAL
*/
func NewDepositoryDAL(c *DepositoryDALConfig) model.DepositoryDAL {
	ddal := &DepositoryDAL{}

	ddal.DepositoryList = make(map[string]*model.Depository)
	ddal.DB = c.DB

	return ddal
}

/*
@func: CreateDepository
@description:

	insert a new depository to depository list
*/
func (d *DepositoryDAL) CreateDepository(ctx context.Context, depository *model.Depository) {
	d.DepositoryList[depository.ClientID] = depository
}

/*
@func: DeleteDepository
@description:

	delete the specified depository from depository list
*/
func (d *DepositoryDAL) DeleteDepository(ctx context.Context, depositoryID string) {
	delete(d.DepositoryList, depositoryID)
}

// CreateDepositoryInRDS create depository core info to rds
func (d *DepositoryDAL) CreateDepositoryInRDS(ctx context.Context, info *model.DepositoryCore) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_cores").Create(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"info":  info,
		}).Warn("Failed to create depository core info to rds")
		return err
	}

	return nil
}

// DeleteDepositoryInRDSByID delete depository core info by id in rds
func (d *DepositoryDAL) DeleteDepositoryInRDSByID(ctx context.Context, id string) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_cores").Where("id=?", id).Delete(&model.DepositoryCore{}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to delete depository core info by id in rds")
		return err
	}
	return nil
}

// GetDepositoryInRDS obtain all depository core info from rds
func (d *DepositoryDAL) GetDepositoryInRDS(ctx context.Context) ([]model.DepositoryCore, error) {
	var infos []model.DepositoryCore

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_cores").Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all depository core info from rds")
		return nil, err
	}

	return infos, nil
}

// GetDepositoryInRDSByID get depository core info by id from rds
func (d *DepositoryDAL) GetDepositoryInRDSByID(ctx context.Context, id string) (*model.DepositoryCore, error) {
	var info model.DepositoryCore
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_cores").Where("id = ?", id).First(&info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to get depository core info by id from rds")
		return nil, err
	}
	return &info, nil
}

// UpdateDepositoryInRDSByID update depository core info by id in rds
func (d *DepositoryDAL) UpdateDepositoryInRDSByID(ctx context.Context, info *model.DepositoryCore) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_cores").Where("id=?", info.ID).Updates(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    info.ID,
		}).Warn("Failed to update depository core info by id in rds")
		return err
	}
	return nil
}

// GetDepositoryBetweenIDInRDS get depository core info Between id from rds
func (d *DepositoryDAL) GetDepositoryBetweenIDInRDS(ctx context.Context, ids []string) ([]model.DepositoryCore, error) {
	var infos []model.DepositoryCore

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_cores").Where("id IN (?)", ids).Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all depository core info from rds")
		return nil, err
	}

	return infos, nil
}

// Clear delete all
func (d *DepositoryDAL) Clear() {
	if err := d.DB.Exec("DELETE FROM public.depository_cores").Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Fail to clear depository core table")
	}
}
