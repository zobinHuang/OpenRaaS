package handler

import (
	"serverd/handler/middleware"
	"serverd/model"
	"serverd/model/apperrors"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/*
	struct: Handler
	description: handler layer
*/
type Handler struct {
	InstanceService model.InstanceService
	StreamerService model.StreamerService
	StreamerClient  *model.Streamer
}

/*
	struct: Config
	description: used for config instance of struct Handler
*/
type Config struct {
	R               *gin.Engine
	InstanceService model.InstanceService
	StreamerService model.StreamerService
	BaseURL         string
	TimeoutDuration time.Duration
}

/*
	func: NewHandler
	description: define endpoints for handler, and map each endpoint to handler func
*/

func NewHandler(c *Config) {
	h := &Handler{
		InstanceService: c.InstanceService,
		StreamerService: c.StreamerService,
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
	// disable authentication for debug
	// g.GET("/wsconnect", middleware.AuthUser(h.TokenService), h.WSConnect)
	g.GET("/wsconnect", h.WSConnect)
	g.POST("/createinstance", h.CreateInstance)
	g.POST("/deleteinstance", h.DeleteInstance)
}
