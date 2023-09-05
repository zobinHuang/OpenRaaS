package model

import (
	"time"

	"gorm.io/gorm"
)

/*
@model: FileStore
@description: filestore client
*/
type FileStore struct {
	FileStoreCore
	Client
}

/*
@model: FileStoreCore
@description: metadata for filestore client
@param SupportApp: slice to json string
*/
type FileStoreCore struct {
	CreateAt              time.Time      `json:"create_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeleteAt              gorm.DeletedAt `gorm:"index" json:"delete_at"`
	ID                    string         `gorm:"unique,not null" json:"id"`
	IP                    string         `gorm:"not null" json:"ip"`
	Port                  int            `gorm:"not null" json:"port"`
	Protocol              string         `gorm:"not null" json:"protocol"`
	Directory             string         `gorm:"not null" json:"directory"`
	Username              string         `json:"username"`
	Password              string         `json:"password"`
	Mem                   float64        `json:"mem"`
	IsContainFastNetspeed bool           `gorm:"not null" json:"is_contain_fast_netspeed"`
}
