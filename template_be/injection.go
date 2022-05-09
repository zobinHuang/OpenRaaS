package main

import (
	"business/dal"
	"business/handler"
	"business/service"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

/*
	func: inject
	description: build layer architecture
*/
func inject(ds *dal.DataSource) (*gin.Engine, error) {
	log.Info("Injecting data sources")

	// --------------------- DAL Layer --------------------------
	rdbDAL := dal.NewRDbDAL(ds.DB)

	// --------------------- Service Layer --------------------------
	rdbService := service.NewRDbService(&service.RDbServiceConfig{
		RDbDAL: rdbDAL,
	})

	// --------------------- Handler Layer --------------------------
	// initialize gin router
	router := gin.Default()

	// obtain base url
	baseURL := os.Getenv("TEST_API_URL")

	// handler timeout
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	// Handler
	handler.NewHandler(&handler.Config{
		R:               router,
		RDbService:      rdbService,
		BaseURL:         baseURL,
		TimeoutDuration: time.Duration(time.Duration(ht) * time.Second),
	})

	return router, nil
}
