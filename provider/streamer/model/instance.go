package model

type Instance struct {
	InstanceCore
}

/*
	@model: StreamInstance
	@description:
		represent an instance of stream application
*/
type StreamInstance struct {
	*StreamApplication
	InstanceID string `json:"instance_id"`
}

type InstanceCore struct {
}
