package model

import (
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
	Port                  int            `gorm:"not null" json:"port"`
	Tag                   string         `json:"tag"`
	Mem                   float64        `json:"mem"`
	IsContainFastNetspeed bool           `gorm:"not null" json:"is_contain_fast_netspeed"`
}
