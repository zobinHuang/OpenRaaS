package dal

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
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

	insert a new depositary to depositary list
*/
func (d *DepositoryDAL) CreateDepository(ctx context.Context, depositary *model.Depository) {
	d.DepositoryList[depositary.ClientID] = depositary
}

/*
@func: DeleteDepository
@description:

	delete the specified depositary from depositary list
*/
func (d *DepositoryDAL) DeleteDepository(ctx context.Context, depositaryID string) {
	delete(d.DepositoryList, depositaryID)
}

// CreateDepositoryInRDS create depositary core info to rds
func (d *DepositoryDAL) CreateDepositoryInRDS(ctx context.Context, info *model.DepositoryCore) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depositary_cores").Create(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"info":  info,
		}).Warn("Failed to create depositary core info to rds")
		return err
	}

	return nil
}

// DeleteDepositoryInRDSByID delete depositary core info by id in rds
func (d *DepositoryDAL) DeleteDepositoryInRDSByID(ctx context.Context, id string) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depositary_cores").Where("id=?", id).Delete(&model.DepositoryCore{}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to delete depositary core info by id in rds")
		return err
	}
	return nil
}

// GetDepositoryInRDS obtain all depositary core info from rds
func (d *DepositoryDAL) GetDepositoryInRDS(ctx context.Context) ([]model.DepositoryCore, error) {
	var infos []model.DepositoryCore

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depositary_cores").Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all depositary core info from rds")
		return nil, err
	}

	return infos, nil
}

// GetDepositoryInRDSByID get depositary core info by id from rds
func (d *DepositoryDAL) GetDepositoryInRDSByID(ctx context.Context, id string) (*model.DepositoryCore, error) {
	var info model.DepositoryCore
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depositary_cores").Where("id = ?", id).First(&info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to get depositary core info by id from rds")
		return nil, err
	}
	return &info, nil
}

// UpdateDepositoryInRDSByID update depositary core info by id in rds
func (d *DepositoryDAL) UpdateDepositoryInRDSByID(ctx context.Context, info *model.DepositoryCore) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depositary_cores").Where("id=?", info.ID).Updates(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    info.ID,
		}).Warn("Failed to update depositary core info by id in rds")
		return err
	}
	return nil
}

// GetDepositoryBetweenIDInRDS get depositary core info Between id from rds
func (d *DepositoryDAL) GetDepositoryBetweenIDInRDS(ctx context.Context, ids []string) ([]model.DepositoryCore, error) {
	var infos []model.DepositoryCore

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depositary_cores").Where("id IN (?)", ids).Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all depositary core info from rds")
		return nil, err
	}

	return infos, nil
}

// Clear delete all
func (d *DepositoryDAL) Clear() {
	d.DB.Delete(&model.DepositoryCore{})
}
