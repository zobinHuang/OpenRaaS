package dal

import (
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
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
	err = db.AutoMigrate(&model.ProviderCore{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.FileStoreCore{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.DepositoryCore{})
	if err != nil {
		return err
	}
	return nil
}
