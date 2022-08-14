package model

import (
	"fmt"

	"github.com/Unknwon/goconfig"
)

type Config struct {
	SchedulerHost   string
	SchedulerPort   string
	SchedulerScheme string
	SchedulerPath   string
}

func (c *Config) Load(group string, key string) string {
	cfg, err := goconfig.LoadConfigFile("config")
	if err != nil {
		fmt.Printf("Unable to load config file.")
		return ""
	}
	var ret string
	ret, err = cfg.GetValue(group, key)
	if err != nil {
		fmt.Printf("Unable to find target %s.%s in config file.", group, key)
		return ""
	}
	return ret
}

func (c *Config) LoadConfigFile() {
	c.SchedulerHost = c.Load("Scheduler", "host")
	c.SchedulerPort = c.Load("Scheduler", "port")
	c.SchedulerPath = c.Load("Scheduler", "path")
	c.SchedulerScheme = c.Load("Scheduler", "scheme")
}
