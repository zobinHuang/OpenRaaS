package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

/*
@model: Provider
@description: provider client
*/
type Provider struct {
	ProviderCore
	Client
}

/*
@model: ProviderCore
@description: metadata for provider client
*/
type ProviderCore struct {
	CreateAt     time.Time      `json:"create_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeleteAt     gorm.DeletedAt `gorm:"index" json:"delete_at"`
	ID           string         `gorm:"unique,not null" json:"id"`
	IP           string         `gorm:"unique,not null" json:"ip"`
	Port         string         `json:"port"`
	Processor    float64        `json:"processor"`
	Bandwidth    float64        `json:bandwidth`
	Latency      float64        `json:latency`
	IsContainGPU bool           `gorm:"not null" json:"is_contain_gpu"`
}

// ProviderCoreWithInst ProviderCore with instance history in blockchain
type ProviderCoreWithInst struct {
	ProviderCore
	InstHistory map[string]string `json:"inst_history"`
}

//////// util ////////

func (s ProviderCoreWithInst) GetMeanHistory() string {
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

func (s ProviderCoreWithInst) GetAbnormalHistoryTimes() int {
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

func (s ProviderCore) DeviceType() string {
	if s.IsContainGPU {
		return "高性能计算设备"
	} else {
		return "普通计算设备"
	}
}

//////// print ////////

// func (s ProviderCoreWithInst) String() string {
// 	// Customize fmt.Println(s)
// 	return s.ID
// }

func (s ProviderCoreWithInst) DetailedInfo() string {
	// Customize fmt.Println(s)
	l1 := fmt.Sprintf("服务提供节点: %s", s.ID)
	l2 := fmt.Sprintf("IP 地址: %s | 服务端口: 8080 | 设备类型: %s", s.IP, s.DeviceType())
	l3 := fmt.Sprintf("计算资源 (GF): %.2f | 网络带宽 (Mbps): %.2f | 设备延迟 (ms): %.2f", s.Processor, s.Bandwidth, s.Latency)

	abt := s.GetAbnormalHistoryTimes()
	var abtStr string
	if abt > 0 {
		abtStr = "属于异常节点"
	} else {
		abtStr = "工作正常"
	}

	l4 := fmt.Sprintf("服务历史: 平均服务延迟为 %s, 存在 %d 次异常行为, %s", s.GetMeanHistory(), s.GetAbnormalHistoryTimes(), abtStr)

	ans := l1 + "\n" + l2 + "\n" + l3 + "\n" + l4 + "\n"

	return ans
}
