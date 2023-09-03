package model

/*
@model: StreamInstance
@description:

	represent an instance of stream application
*/
type StreamInstance struct {
	*StreamApplication
	InstanceID   string `json:"instance_id"`
	ScreenWidth  int    `json:"screen_width"`
	ScreenHeight int    `json:"screen_height"`
	VCodec       string `json:"vcodec"`
	FPS          int    `json:"fps"`
}

/*
@model: StreamInstanceRoom
@description:

	room of a initialized stream application instance
*/
type StreamInstanceRoom struct {
	*StreamInstance

	Provider     *Provider
	ConsumerList map[string]*Consumer

	SelectedDepository      *Depository
	PotentialDepositoryList map[string]*Depository

	SelectedFileStore      *FileStore
	PotentialFileStoreList map[string]*FileStore
}
