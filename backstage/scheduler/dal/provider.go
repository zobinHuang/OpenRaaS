package dal

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
@struct: ProviderDAL
@description: DAL layer
*/
type ProviderDAL struct {
	DB           *gorm.DB
	ProviderList map[string]*model.Provider
}

/*
@struct: ProviderDALConfig
@description: used for config instance of struct ProviderDAL
*/
type ProviderDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewProviderDAL
@description:

	create, config and return an instance of struct ProviderDAL
*/
func NewProviderDAL(c *ProviderDALConfig) model.ProviderDAL {
	pdal := &ProviderDAL{}

	pdal.ProviderList = make(map[string]*model.Provider)
	pdal.DB = c.DB

	return pdal
}

/*
@func: CreateProvider
@description:

	insert a new provider to provider list
*/
func (d *ProviderDAL) CreateProvider(ctx context.Context, provider *model.Provider) {
	d.ProviderList[provider.ClientID] = provider
}

/*
@func: DeleteProvider
@description:

	delete the specified provider from provider list
*/
func (d *ProviderDAL) DeleteProvider(ctx context.Context, providerID string) {
	delete(d.ProviderList, providerID)
}

func (d *ProviderDAL) GetProvider() []*model.Provider {
	providers := make([]*model.Provider, 0, 0)
	for _, value := range d.ProviderList {
		providers = append(providers, value)
	}
	return providers
}

// CreateProviderInRDS create provider core info to rds
func (d *ProviderDAL) CreateProviderInRDS(ctx context.Context, provider *model.ProviderCore) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_cores").Create(provider).Error; err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"provider": provider,
		}).Warn("Failed to create provider core info to rds")
		return err
	}

	return nil
}

// DeleteProviderInRDSByID
func (d *ProviderDAL) DeleteProviderInRDSByID(ctx context.Context, id string) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_cores").Where("id=?", id).Delete(&model.ProviderCore{}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to delete provider core info by id in rds")
		return err
	}

	return nil
}

// GetProviderInRDS obtain all provider core info from rds
func (d *ProviderDAL) GetProviderInRDS(ctx context.Context) ([]model.ProviderCore, error) {
	var infos []model.ProviderCore

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_cores").Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all provider core info from rds")
		return nil, err
	}

	return infos, nil
}

// GetProviderInRDSByID get provider core info by id from rds
func (d *ProviderDAL) GetProviderInRDSByID(ctx context.Context, id string) (*model.ProviderCore, error) {
	var info model.ProviderCore
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_cores").Where("id = ?", id).First(&info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to get provider core info by id from rds")
		return nil, err
	}
	return &info, nil
}

// UpdateProviderInRDSByID update provider core info by id in rds
func (d *ProviderDAL) UpdateProviderInRDSByID(ctx context.Context, provider *model.ProviderCore) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_cores").Where("id=?", provider.ID).Updates(provider).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    provider.ID,
		}).Warn("Failed to update provider core info by id in rds")
		return err
	}
	return nil
}

// Clear delete all
func (d *ProviderDAL) Clear() {
	if err := d.DB.Exec("DELETE FROM public.provider_cores").Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Fail to clear provider_cores table")
	}
}
