package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/handler/middleware"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model/apperrors"
)

/*
	struct: Handler
	description: handler layer
*/
type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
}

/*
	struct: Config
	description: used for config instance of struct Handler
*/
type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	BaseURL         string
	TimeoutDuration time.Duration
}

/*
	func: NewHandler
	description: define endpoints for handler, and map each endpoint to handler func
*/

func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
	}

	// response to cors request (accept all origins)
	c.R.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))

	// group to base url
	g := c.R.Group(c.BaseURL)

	// add timeout middleware
	g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))

	// user authetication
	g.POST("/signin", h.Signin)
	g.POST("/signup", h.Signup)
	g.POST("/signout", middleware.AuthUser(h.TokenService), h.SignOut)
}
