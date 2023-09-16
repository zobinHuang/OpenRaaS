package model

import (
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
	IsContainGPU bool           `gorm:"not null" json:"is_contain_gpu"`
}

// ProviderCoreWithInst ProviderCore with instance history in blockchain
type ProviderCoreWithInst struct {
	ProviderCore
	InstHistory map[string]string `json:"inst_history"`
}
