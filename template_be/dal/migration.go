package dal

import (
	"business/model"
	"gorm.io/gorm"
)

/*
	func: dBMigrator
	description: migrate models to database tables
*/
func dBMigrator(db *gorm.DB) error {
	err := db.AutoMigrate(&model.RDbModel{})
	if err != nil {
		return err
	}
	return nil
}
