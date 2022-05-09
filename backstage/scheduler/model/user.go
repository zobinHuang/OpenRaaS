package model

import (
	"time"

	"gorm.io/gorm"
)

/*
	model: Test
	description: relation-database model
*/
type User struct {
	CreateAt  time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeleteAt  gorm.DeletedAt `gorm:"index" json:"delete_at"`
	Id        uint64         `gorm:"AUTO_INCREMENT" json:"uid"`
	Email     string         `gorm:"unique,not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
}

/*
	func: TableName
	description: realize Tabler interface of GORM
*/
func (u *User) TableName() string {
	return "user"
}
