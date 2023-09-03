package model

import (
	"gorm.io/gorm"
	"time"
)

/*
@model: Depositary
@description: depositary client
*/
type Depositary struct {
	DepositaryCore
	Client
}

/*
@model: DepositaryCore
@description: metadata for depositary client
*/
type DepositaryCore struct {
	CreateAt              time.Time      `json:"create_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeleteAt              gorm.DeletedAt `gorm:"index" json:"delete_at"`
	ID                    string         `gorm:"unique,not null" json:"id"`
	IP                    string         `gorm:"not null" json:"ip"`
	Port                  int            `gorm:"not null" json:"port"`
	Tag                   string         `json:"tag"`
	IsContainFastNetspeed bool           `gorm:"not null" json:"is_contain_fast_netspeed"`
}
