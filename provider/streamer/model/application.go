package model

import (
	"time"

	"gorm.io/gorm"
)

/*
	@model: StreamApplication
	@description:
		represent an stream application
*/
type StreamApplication struct {
	StreamApplicationCore
	ApplicationCore
}

/*
	@model: StreamApplicationCore
	@description:
		metadata for stream application
*/
type StreamApplicationCore struct {
}

/*
	@model: Application
	@description:
		common meta data of an application
*/
type ApplicationCore struct {
	CreateAt        time.Time      `json:"create_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeleteAt        gorm.DeletedAt `gorm:"index" json:"delete_at"`
	ApplicationName string         `gorm:"not null" json:"application_name"`
	ApplicationID   string         `gorm:"unique,not null" json:"application_id"`
	ApplicationPath string         `json:"application_path"`
	ApplicationFile string         `json:"applictaion_file"`
	HWKey           string         `json:"hwkey"`
	OperatingSystem string         `gorm:"not null" json:"operating_system"`
	CreateUser      string         `gorm:"not null" json:"create_user"`
	Description     string         `json:"description"`
	UsageCount      int64          `json:"usage_count"`
}
