package handler

import (
	"business/handler/middleware"
	"business/model"
	"business/model/apperrors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

/*
	struct: Handler
	description: handler layer
*/
type Handler struct {
	RDbService model.RDbService
}

/*
	struct: Config
	description: used for config instance of struct Handler
*/
type Config struct {
	R               *gin.Engine
	RDbService      model.RDbService
	BaseURL         string
	TimeoutDuration time.Duration
}

/*
	func: NewHandler
	description: define endpoints for handler, and map each endpoint to handler func
*/

func NewHandler(c *Config) {
	h := &Handler{
		RDbService: c.RDbService,
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
	g.POST("/test", h.Test)
}
