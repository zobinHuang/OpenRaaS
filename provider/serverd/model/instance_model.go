package model

/*
	model: InstanceModel
	description: store attributes of a instance (inside a container)
*/
type InstanceModel struct {
	InstanceCore
	InstanceConnection
	Done             chan (struct{}) `json:"-"` // a channel used for close the vm
	AppPath          string          `json:"application_path"`
	AppFile          string          `json:"application_file"`
	AppName          string          `json:"application_name"`
	HWKey            string          `json:"hw_key"` // 'app' or 'game'
	ScreenWidth      int             `json:"screen_width"`
	ScreenHeight     int             `json:"screen_height"`
	AppOption        string          `json:"app_option"`
	FPS              int             `json:"fps"`
	VCodec           string          `json:"vcodec"`
	ImageName        string          `json:"image_name"`
	FilestoreList    []Filestore     `json:"filestore_list"`  // list of ip address
	DepositoryList   []Depository    `json:"depository_list"` // list of ip address
	TargetFilestore  Filestore       `json:"target_filestore"`
	TargetDepository Depository      `json:"target_depository"`
	RunInLinux       bool            `json:"run_in_linux"`
	RunWithGpu       bool
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
	@model: Depository
	@description: depository client
*/
type Depository struct {
	DepositoryCore
	InstHistory map[string]string `json:"inst_history"`
}

/*
	@model: DepositoryCore
	@description: metadata for depository client
	@example:
		depository.HostAddress = "127.0.0.1"
		depository.Port = "5000"
		depository.Tag = "latest"
*/
type DepositoryCore struct {
	HostAddress           string `json:"IP"`
	Port                  string `json:"port"`
	Tag                   string `json:"tag"`
	IsContainFastNetspeed bool   `json:"is_contain_fast_netspeed"`
}

/*
	@model: Filestore
	@description: filestore client
*/
type Filestore struct {
	FilestoreCore
	InstHistory map[string]string `json:"inst_history"`
}

/*
	@model: FilestoreCore
	@description: metadata for filestore client
	@example:
		HostAddress "192.168.10.189"
		Port        "7189"
		Protocol    "davfs"
		Username    "kb109"
		Password    "******"
		Directory   "/public_hdd/game/PC/dcwine"
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
