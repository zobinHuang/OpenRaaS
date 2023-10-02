package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

/*
@model: Depository
@description: depository client
*/
type Depository struct {
	DepositoryCore
	Client
}

/*
@model: DepositoryCore
@description: metadata for depository client
@param SupportApp: slice to json string
*/
type DepositoryCore struct {
	CreateAt              time.Time      `json:"create_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeleteAt              gorm.DeletedAt `gorm:"index" json:"delete_at"`
	ID                    string         `gorm:"unique,not null" json:"id"`
	IP                    string         `gorm:"not null" json:"ip"`
	Port                  string         `gorm:"not null" json:"port"`
	Tag                   string         `json:"tag"`
	Mem                   float64        `json:"mem"`
	Bandwidth             float64        `json:bandwidth`
	Latency               float64        `json:latency`
	IsContainFastNetspeed bool           `gorm:"not null" json:"is_contain_fast_netspeed"`
}

// DepositoryCoreWithInst DepositoryCore with instance history in blockchain
type DepositoryCoreWithInst struct {
	DepositoryCore
	InstHistory map[string]string `json:"inst_history"`
}

//////// util ////////

func (s DepositoryCoreWithInst) GetMeanHistory() string {
	// 获取服务历史的均值
	total := 0.
	count := 0
	unit := "bitsPerSecond"

	for _, value := range s.InstHistory {
		strValue := strings.TrimSuffix(value, unit)
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			total += floatValue
			count++
		}
	}

	newUnit := "mbps"

	if count > 0 {
		average := float64(total) / float64(count) / 1e6 // bps 转成 mbps
		return fmt.Sprintf("%.2f %s", average, newUnit)
	} else {
		fmt.Println("字典中没有有效的值.")
		return ""
	}
}

func (s DepositoryCoreWithInst) GetAbnormalHistoryTimes() int {
	// 从服务历史获取异常行为的次数
	count := 0
	unit := "bitsPerSecond"

	for _, value := range s.InstHistory {
		strValue := strings.TrimSuffix(value, unit)
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			if floatValue < 30*1e6 {
				count++
			}
		}
	}

	return count
}

func (s DepositoryCore) DeviceType() string {
	if s.IsContainFastNetspeed {
		return "高性能存储设备"
	} else {
		return "普通存储设备"
	}
}

//////// print ////////

// func (s DepositoryCoreWithInst) String() string {
// 	// Customize fmt.Println(s)
// 	return s.ID
// }

func (s DepositoryCoreWithInst) DetailedInfo() string {
	// Customize fmt.Println(s)
	l1 := fmt.Sprintf("镜像仓库节点: %s", s.ID)
	l2 := fmt.Sprintf("IP 地址: %s | 服务端口: %s | 设备类型: %s", s.IP, s.Port, s.DeviceType())
	l3 := fmt.Sprintf("存储资源 (GB): %.2f | 网络带宽 (Mbps): %.2f | 设备延迟 (ms): %.2f", s.Mem, s.Bandwidth, s.Latency)
	l4 := fmt.Sprintf("仓库类型: Docker Registry | 镜像前缀: '%s:%s/' | 镜像后缀: ':%s'", s.IP, s.Port, s.Tag)

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
