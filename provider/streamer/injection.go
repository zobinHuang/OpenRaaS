package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/zobinHuang/OpenRaaS/provider/streamer/dal"
	"github.com/zobinHuang/OpenRaaS/provider/streamer/handler"
	"github.com/zobinHuang/OpenRaaS/provider/streamer/service"
	"github.com/zobinHuang/OpenRaaS/provider/streamer/service/communicator"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

/*
	func: inject
	description: build layer architecture
*/
func inject() (*gin.Engine, error) {
	log.Info("Injecting data sources")

	// --------------------- DAL Layer --------------------------
	instanceDAL := dal.NewInstanceDAL(&dal.InstanceDALConfig{})
	schedulerDAL := dal.NewSchedulerDAL(&dal.SchedulerDALConfig{})
	daemonDAL := dal.NewDaemonDAL(&dal.DaemonDALConfig{})
	webRTCStreamDAL := dal.NewWebRTCStreamDAL(&dal.WebRTCStreamDALConfig{})

	// --------------------- Comnunicator Layer ---------------------
	wsCommunicator := communicator.NewWebsocketCommunicator(&communicator.WebsocketCommunicatorConfig{
		SchedulerDAL:    schedulerDAL,
		DaemonDAL:       daemonDAL,
		InstanceDAL:     instanceDAL,
		WebRTCStreamDAL: webRTCStreamDAL,
	})

	// --------------------- Service Layer --------------------------
	instanceService := service.NewInstanceService(&service.InstanceServiceConfig{
		InstanceDAL: instanceDAL,
	})

	schedulerService := service.NewSchedulerService(&service.SchedulerServiceConfig{
		WebsocketCommunicator: wsCommunicator,
		SchedulerDAL:          schedulerDAL,
	})

	daemonService := service.NewDaemonService(&service.DaemonServiceConfig{
		WebsocketCommunicator: wsCommunicator,
		DaemonDAL:             daemonDAL,
	})

	// --------------------- Handler Layer --------------------------
	// initialize gin router
	router := gin.Default()

	// obtain base url
	baseURL := os.Getenv("STREAMER_API_URL")

	// handler timeout
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	// Handler
	handler.NewHandler(&handler.Config{
		R:                     router,
		BaseURL:               baseURL,
		InstanceService:       instanceService,
		SchedulerService:      schedulerService,
		DaemonService:         daemonService,
		WebsocketCommunicator: wsCommunicator,
		TimeoutDuration:       time.Duration(time.Duration(ht) * time.Second),
	})

	return router, nil
}
