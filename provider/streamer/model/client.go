package model

/*
	@model: DepositoryCore
	@description: metadata for depository client
*/
type DepositoryCore struct {
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
