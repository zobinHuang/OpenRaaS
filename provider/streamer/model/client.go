package model

/*
	@model: DepositaryCore
	@description: metadata for depositary client
*/
type DepositoryCore struct {
	HostAddress string `json:"IP"`
	Port        string `json:"port"`
	Tag         string `json:"tag"`
}

/*
	@model: FilestoreCore
	@description: metadata for filestore client
*/
type FilestoreCore struct {
	HostAddress string `json:"IP"`
	Port        string `json:"port"`
	Protocol    string `json:"protocol"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Directory   string `json:"directory"`
}
