package model

import (
	"gorm.io/gorm"
	"time"
)

/*
	model: Test
	description: relation-database model
*/
type RDbModel struct {
	CreateAt  time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeleteAt  gorm.DeletedAt `gorm:"index" json:"delete_at"`
	Id        uint64         `gorm:"AUTO_INCREMENT" json:"uid"`
	UserName  string         `gorm:"unique,not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
}

/*
	func: TableName
	description: realize Tabler interface of GORM
*/
func (uc *RDbModel) TableName() string {
	return "test_model"
}
