package dal

import (
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
	"gorm.io/gorm"
)

/*
@func: dBMigrator
@description: migrate models to database tables
*/
func dBMigrator(db *gorm.DB) error {
	// migrate stream application
	err := db.AutoMigrate(&model.StreamApplication{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.ProviderCoreWithInst{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.FileStoreCoreWithInst{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.DepositoryCoreWithInst{})
	if err != nil {
		return err
	}
	return nil
}
