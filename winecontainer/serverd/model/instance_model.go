package model

/*
	model: InstanceModel
	description: store attributes of a instance (inside a container)
*/
type InstanceModel struct {
	InstanceCore
	InstanceConnection
	Done             chan (struct{})  `json:"-"` // a channel used for close the vm
	AppPath          string           `json:"application_path"`
	AppFile          string           `json:"application_file"`
	AppName          string           `json:"application_name"`
	HWKey            string           `json:"hw_key"` // 'app' or 'game'
	ScreenWidth      int              `json:"screen_width"`
	ScreenHeight     int              `json:"screen_height"`
	WineOption       string           `json:"wine_option"`
	FPS              int              `json:"fps"`
	VCodec           string           `json:"vcodec"`
	FilestoreList    []FilestoreCore  `json:"filestore_list"`  // list of ip address
	DepositaryList   []DepositaryCore `json:"depositary_list"` // list of ip address
	TargetFilestore  FilestoreCore    `json:"target_filestore"`
	TargetDepositary DepositaryCore   `json:"target_depositary"`
}

/*
	model: InstanceConnection
	description: store attributes of a connection between wine container (instance) and streamer
*/
type InstanceConnection struct {
	InstanceIP   string `json:"instance_ip"`
	VideoRTCPort string `json:"video_rtc_port"`
	AudioRTCPort string `json:"audio_rtc_port"`
	InputPort    string `json:"input_port"`
}

/*
	model: DeleteInstanceModel
	description: store attributes of a instance which is ordered to shut down
*/
type DeleteInstanceModel struct {
	InstanceCore
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
	@model: Depositary
	@description: depositary client
*/
type Depositary struct {
	DepositaryCore
}

/*
	@model: DepositaryCore
	@description: metadata for depositary client
*/
type DepositaryCore struct {
	HostAddress string `json:"host_address"`
	Port        string `json:"port"`
	Tag         string `json:"tag"`
}

/*
	@model: Filestore
	@description: filestore client
*/
type Filestore struct {
	FilestoreCore
}

/*
	@model: FilestoreCore
	@description: metadata for filestore client
	@example:
		HostAddress "192.168.10.189"
		Port        "7189"
		Protocal    "davfs"
		Username    "kb109"
		Password    "******"
		Directory   "/public_hdd/game/PC/dcwine"
*/
type FilestoreCore struct {
	HostAddress string `json:"host_address"`
	Port        string `json:"port"`
	Protocal    string `json:"protocal"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Directory   string `json:"directory"`
}
