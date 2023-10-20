package model

import (
	"fmt"
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

func (s StreamApplication) ProviderReq() string {
	if s.IsProviderReqGPU {
		return "高性能 (GPU)"
	} else {
		return "普通"
	}
}

func (s StreamApplication) FilestoreReq() string {
	if s.IsFileStoreReqFastNetspeed {
		return "高速读写"
	} else {
		return "普通"
	}
}

func (s StreamApplication) DepositoryReq() string {
	if s.IsDepositoryReqFastNetspeed {
		return "高速读取"
	} else {
		return "普通"
	}
}

func (s StreamApplication) DetailedInfo() string {
	// Customize fmt.Println(s)
	l1 := fmt.Sprintf("软件ID: %s | 软件名: %s ", s.ApplicationID, s.ApplicationName)
	l2 := fmt.Sprintf("软件路径: %s | 启动文件: %s | 类型: %s | 镜像名: %s", s.ApplicationPath, s.ApplicationFile, s.HWKey, s.ImageName)
	l3 := fmt.Sprintf("部署的节点: %s | 软件说明: %s", s.FileStoreList, s.Description)
	l4 := fmt.Sprintf("计算资源需求: %s", s.ProviderReq())
	l5 := fmt.Sprintf("存储资源需求 (读写): %s", s.FilestoreReq())
	l6 := fmt.Sprintf("存储资源需求 (只读): %s", s.DepositoryReq())

	ans := l1 + "\n" + l2 + "\n" + l3 + "\n" + l4 + "\n" + l5 + "\n" + l6 + "\n"

	return ans
}
