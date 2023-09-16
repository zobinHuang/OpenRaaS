package model

/*
	@model: DepositoryCore
	@description: metadata for depository client
*/
type DepositoryCore struct {
	HostAddress           string `json:"IP"`
	Port                  string `json:"port"`
	Tag                   string `json:"tag"`
	IsContainFastNetspeed bool   `json:"is_contain_fast_netspeed"`
}

type DepositoryCoreWithInst struct {
	DepositoryCore
	InstHistory map[string]string `json:"inst_history"`
}

/*
	@model: FilestoreCore
	@description: metadata for filestore client
*/
type FilestoreCore struct {
	HostAddress           string `json:"IP"`
	Port                  string `json:"port"`
	Protocol              string `json:"protocol"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	Directory             string `json:"directory"`
	IsContainFastNetspeed bool   `json:"is_contain_fast_netspeed"`
}

type FilestoreCoreWithInst struct {
	FilestoreCore
	InstHistory map[string]string `json:"inst_history"`
}
