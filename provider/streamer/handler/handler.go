package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zobinHuang/BrosCloud/provider/streamer/handler/middleware"
	"github.com/zobinHuang/BrosCloud/provider/streamer/model"
	"github.com/zobinHuang/BrosCloud/provider/streamer/model/apperrors"
)

/*
	struct: Handler
	description: handler layer
*/
type Handler struct {
	InstanceService       model.InstanceService
	SchedulerService      model.SchedulerService
	DaemonService         model.DaemonService
	WebsocketCommunicator model.WebsocketCommunicator
}

/*
	struct: Config
	description: used for config instance of struct Handler
*/
type Config struct {
	R                     *gin.Engine
	InstanceService       model.InstanceService
	SchedulerService      model.SchedulerService
	DaemonService         model.DaemonService
	WebsocketCommunicator model.WebsocketCommunicator
	BaseURL               string
	TimeoutDuration       time.Duration
}

/*
	func: NewHandler
	description: define endpoints for handler, and map each endpoint to handler func
*/

func NewHandler(c *Config) {
	h := &Handler{
		InstanceService:       c.InstanceService,
		SchedulerService:      c.SchedulerService,
		DaemonService:         c.DaemonService,
		WebsocketCommunicator: c.WebsocketCommunicator,
	}

	// response to cors request (accept all origins)
	c.R.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// group to base url
	g := c.R.Group(c.BaseURL)

	// add timeout middleware
	g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))

	g.GET("/wsconnect", h.WebsocketConnect)
}
