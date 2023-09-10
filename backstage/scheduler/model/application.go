package model

import (
	"time"

	"gorm.io/gorm"
)

/*
@enum
@description: application type
*/
const (
	APPLICATIOON_TYPE_STREAM  string = "stream"
	APPLICATIOON_TYPE_CONSOLE string = "console"
)

/*
@enum
@description: order scheme (used by database searching)
*/
const (
	ORDER_BY_UPDATE_TIME string = "orderByUpdateTime"
	ORDER_BY_NAME        string = "orderByName"
	ORDER_BY_USAGE_COUNT string = "orderByUsageCount"
)

/*
@model: StreamApplication
@description:

	represent an stream application
*/
type StreamApplication struct {
	StreamApplicationCore
	ApplicationCore
	AppInfoAttach
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

/*
@model: Application
@description:

	request info of an application
*/
type AppInfoAttach struct {
	FileStoreList               string `json:"filestore_list"`
	ImageName                   string `json:"image_name"`
	IsProviderReqGPU            bool   `gorm:"not null" json:"is_provider_req_gpu"`
	IsFileStoreReqFastNetspeed  bool   `gorm:"not null" json:"is_filestore_req_fast_netspeed"`
	IsDepositoryReqFastNetspeed bool   `gorm:"not null" json:"is_depository_req_fast_netspeed"`
}
