package model

/*
	model: InstanceModel
	description:
		store attributes of a stream instance (i.e. view from daemon)
*/
type StreamInstanceDaemonModel struct {
	InstanceCore
	InstanceConnection
	AppPath          string           `json:"application_path"`
	AppFile          string           `json:"application_file"`
	AppName          string           `json:"application_name"`
	HWKey            string           `json:"hw_key"`
	ScreenWidth      int              `json:"screen_width"`
	ScreenHeight     int              `json:"screen_height"`
	FPS              int              `json:"fps"`
	VCodec           string           `json:"vcodec"`
	ImageName        string           `json:"image_name"`
	FilestoreList    []FilestoreCore  `json:"filestore_list"`
	DepositoryList   []DepositoryCore `json:"depository_list"`
	TargetFilestore  FilestoreCore    `json:"target_filestore"`
	TargetDepository DepositoryCore   `json:"target_depository"`
}

/*
	model: InstanceConnection
	description:
		store attributes of a connection between wine container (instance) and streamer
*/
type InstanceConnection struct {
	InstanceIP   string `json:"instance_ip"`
	VideoRTCPort string `json:"video_rtc_port"`
	AudioRTCPort string `json:"audio_rtc_port"`
	InputPort    string `json:"input_port"`
}

/*
	model: InstanceCore
	description: metadata for instance
*/
type InstanceCore struct {
	VMID       int    `json:"vmid"`
	Instanceid string `json:"instanceid"`
}

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
