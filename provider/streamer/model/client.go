package model

/*
	@model: DepositaryCore
	@description: metadata for depositary client
*/
type DepositaryCore struct {
	HostAddress string `json:"host_address"`
	Port        string `json:"port"`
}

/*
	@model: FilestoreCore
	@description: metadata for filestore client
*/
type FilestoreCore struct {
	HostAddress string `json:"host_address"`
	Port        string `json:"port"`
}
