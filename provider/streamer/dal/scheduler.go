package dal

import "github.com/zobinHuang/BrosCloud/provider/streamer/model"

type SchedulerDAL struct {
	ICEServers string `json:"iceservers"`
}

type SchedulerDALConfig struct {
}

func NewSchedulerDAL(c *SchedulerDALConfig) model.SchedulerDAL {
	return &SchedulerDAL{}
}

func (d *SchedulerDAL) SetICEServers(iceServer string) {
	d.ICEServers = iceServer
}
