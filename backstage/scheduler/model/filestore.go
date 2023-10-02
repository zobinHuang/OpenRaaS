package model

import (
	"fmt"
	"strconv"
	"strings"
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
	Port                  string         `gorm:"not null" json:"port"`
	Protocol              string         `gorm:"not null" json:"protocol"`
	Directory             string         `gorm:"not null" json:"directory"`
	Username              string         `json:"username"`
	Password              string         `json:"password"`
	Mem                   float64        `json:"mem"`
	Bandwidth             float64        `json:bandwidth`
	Latency               float64        `json:latency`
	IsContainFastNetspeed bool           `gorm:"not null" json:"is_contain_fast_netspeed"`
}

// FileStoreCoreWithInst FileStoreCore with instance history in blockchain
type FileStoreCoreWithInst struct {
	FileStoreCore
	InstHistory map[string]string `json:"inst_history"`
}

//////// util ////////

func (s FileStoreCoreWithInst) GetMeanHistory() string {
	// 获取服务历史的均值
	total := 0.
	count := 0
	unit := "ms"

	for _, value := range s.InstHistory {
		strValue := strings.TrimSuffix(value, unit)
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			total += floatValue
			count++
		}
	}

	if count > 0 {
		average := float64(total) / float64(count)
		return fmt.Sprintf("%.2f %s", average, unit)
	} else {
		fmt.Println("字典中没有有效的值.")
		return ""
	}
}

func (s FileStoreCoreWithInst) GetAbnormalHistoryTimes() int {
	// 从服务历史获取异常行为的次数
	count := 0
	unit := "ms"

	for _, value := range s.InstHistory {
		strValue := strings.TrimSuffix(value, unit)
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			if floatValue > 100 {
				count++
			}
		}
	}

	return count
}

func (s FileStoreCore) DeviceType() string {
	if s.IsContainFastNetspeed {
		return "高性能存储设备"
	} else {
		return "普通存储设备"
	}
}

//////// print ////////

// func (s FileStoreCoreWithInst) String() string {
// 	// Customize fmt.Println(s)
// 	return s.ID
// }

func (s FileStoreCoreWithInst) DetailedInfo() string {
	// Customize fmt.Println(s)
	l1 := fmt.Sprintf("内容存储节点: %s", s.ID)
	l2 := fmt.Sprintf("IP 地址: %s | 服务端口: %s | 设备类型: %s", s.IP, s.Port, s.DeviceType())
	l3 := fmt.Sprintf("存储资源 (GB): %.2f | 网络带宽 (Mbps): %.2f | 设备延迟 (ms): %.2f", s.Mem, s.Bandwidth, s.Latency)
	l4 := fmt.Sprintf("文件传输协议: %s | 共享目录路径: %s | 用户名: %s | 密码: %s", s.Protocol, s.Directory, s.Username, s.Password)

	abt := s.GetAbnormalHistoryTimes()
	var abtStr string
	if abt > 0 {
		abtStr = "属于异常节点"
	} else {
		abtStr = "工作正常"
	}

	l5 := fmt.Sprintf("服务历史: 平均服务延迟为 %s, 存在 %d 次异常行为, %s", s.GetMeanHistory(), s.GetAbnormalHistoryTimes(), abtStr)

	ans := l1 + "\n" + l2 + "\n" + l3 + "\n" + l4 + "\n" + l5 + "\n"

	return ans
}
