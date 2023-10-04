package model

import (
	"fmt"
	"strconv"
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

func (s StreamApplication) DetailedInfo() string {
	// Customize fmt.Println(s)
	l1 := fmt.Sprintf("软件 ID: %s | 软件名: %s ", s.ApplicationID, s.ApplicationName)
	l2 := fmt.Sprintf("软件路径: %s | 启动文件: %s | 软件类型: %s | 镜像 ID: %s", s.ApplicationPath, s.ApplicationFile, s.HWKey, s.ImageName)
	l3 := fmt.Sprintf("支持的内容存储节点: %s | 软件说明: %s", s.FileStoreList, s.Description)
	l4 := fmt.Sprintf("是否需要高性能服务提供节点：%s", strconv.FormatBool(s.IsProviderReqGPU))
	l5 := fmt.Sprintf("是否需要高性能内容存储节点：%s", strconv.FormatBool(s.IsFileStoreReqFastNetspeed))
	l6 := fmt.Sprintf("是否需要高性能镜像仓库节点：%s", strconv.FormatBool(s.IsDepositoryReqFastNetspeed))

	ans := l1 + "\n" + l2 + "\n" + l3 + "\n" + l4 + "\n" + l5 + "\n" + l6 + "\n"

	return ans
}
