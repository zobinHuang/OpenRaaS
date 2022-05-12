package dal

import "github.com/zobinHuang/BrosCloud/provider/streamer/model"

type SchedulerDAL struct {
	ICEServers []string
}

type SchedulerDALConfig struct {
}

func NewSchedulerDAL(c *SchedulerDALConfig) model.SchedulerDAL {
	return &SchedulerDAL{}
}

func (d *SchedulerDAL) AddICEServers(iceServer string) {
	d.ICEServers = append(d.ICEServers, iceServer)
}
