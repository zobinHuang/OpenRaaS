package dal

import "github.com/zobinHuang/BrosCloud/provider/streamer/model"

type DaemonDAL struct {
}

type DaemonDALConfig struct {
}

func NewDaemonDAL(c *DaemonDALConfig) model.DaemonDAL {
	return &DaemonDAL{}
}
