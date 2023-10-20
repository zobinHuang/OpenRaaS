package dal

import (
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model"
	"gorm.io/gorm"
)

/*
	func: dBMigrator
	description: migrate models to database tables
*/
func dBMigrator(db *gorm.DB) error {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}
	return nil
}
