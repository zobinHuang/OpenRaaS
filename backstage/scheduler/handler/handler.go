package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/handler/middleware"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model/apperrors"
)

/*
struct: Handler
description: handler layer
*/
type Handler struct {
	RDbService         model.RDbService
	TokenService       model.TokenService
	ConsumerService    model.ConsumerService
	ProviderService    model.ProviderService
	ApplicationService model.ApplicationService
}

/*
struct: Config
description: used for config instance of struct Handler
*/
type Config struct {
	R                  *gin.Engine
	RDbService         model.RDbService
	TokenService       model.TokenService
	ConsumerService    model.ConsumerService
	ProviderService    model.ProviderService
	ApplicationService model.ApplicationService
	BaseURL            string
	TimeoutDuration    time.Duration
}

/*
func: NewHandler
description: define endpoints for handler, and map each endpoint to handler func
*/
func NewHandler(c *Config) {
	h := &Handler{
		RDbService:         c.RDbService,
		TokenService:       c.TokenService,
		ConsumerService:    c.ConsumerService,
		ProviderService:    c.ProviderService,
		ApplicationService: c.ApplicationService,
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
	// g.GET("/node_online", middleware.AuthUser(h.TokenService), h.NodeOnline)
	g.POST("node_online", h.NodeOnline)

	// g.GET("/application_online", middleware.AuthUser(h.TokenService), h.ApplicationOnline)
	g.POST("application_online", h.ApplicationOnline)

	// disable authentication for debug
	// g.GET("/wsconnect", middleware.AuthUser(h.TokenService), h.WSConnect)
	g.GET("/wsconnect", h.WSConnect)

	// disable authentication for debug
	//g.GET("/application_list", middleware.AuthUser(h.TokenService), h.GetApplictaionList)
	g.GET("/application_list", h.GetApplictaionList)

	// disable authentication for debug
	//g.GET("/application_amount", middleware.AuthUser(h.TokenService), h.GetApplictaionAmount)
	g.GET("/application_amount", h.GetApplictaionAmount)

	// disable authertication for debug
	// g.GET("/application_details", middleware.AuthUser(h.TokenService), h.GetApplictaionDetails)
	g.GET("/application_details", h.GetApplictaionDetails)
}
