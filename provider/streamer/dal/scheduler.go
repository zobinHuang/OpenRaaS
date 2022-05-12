package dal

import "github.com/zobinHuang/BrosCloud/provider/streamer/model"

/*
	@struct: SchedulerDAL
	@description: DAL layer
*/
type SchedulerDAL struct {
	ICEServers []string
}

/*
	@struct: SchedulerDALConfig
	@description: used for config instance of struct SchedulerDAL
*/
type SchedulerDALConfig struct {
}

/*
	@function: NewSchedulerDAL
	@description:
		create, config and return an instance of struct SchedulerDAL
*/
func NewSchedulerDAL(c *SchedulerDALConfig) model.SchedulerDAL {
	return &SchedulerDAL{}
}

/*
	@function: AddICEServers
	@description:
		append new ice servers into global slice
*/
func (d *SchedulerDAL) AddICEServers(iceServer string) {
	d.ICEServers = append(d.ICEServers, iceServer)
}

/*
	@function: GetICEServers
	@description:
		obtain ice servers
*/
func (d *SchedulerDAL) GetICEServers() []string {
	return d.ICEServers
}
