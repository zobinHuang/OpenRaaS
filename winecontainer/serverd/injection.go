package main

import (
	"context"
	"fmt"
	"log"
	"serverd/dal"
	"serverd/handler"
	"serverd/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*
	func: inject
	description: build layer architecture
*/
func inject(ds *dal.DataSource) (*gin.Engine, chan struct{}, error) {
	log.Println("Injecting data sources")

	// --------------------- DAL Layer --------------------------
	// rdbDAL := dal.NewRDbDAL(ds.DB)

	// --------------------- Service Layer --------------------------
	// rdbService := service.NewRDbService(&service.RDbServiceConfig{
	// 	RDbDAL: rdbDAL,
	// })

	fmt.Println("Injection!")

	instanceService := service.NewInstanceService(&service.InstanceServiceConfig{})
	ss := service.NewStreamerService(&service.StreamerServiceConfig{})

	// initiate provider streamer container
	ctx := context.Background()
	ss.RunStreamerContainer(ctx)

	// --------------------- Handler Layer --------------------------
	// initialize gin router
	router := gin.Default()

	// obtain base url ("/api/daemon")
	// baseURL := os.Getenv("SERVERD_API_URL")
	baseURL := "/api/daemon"

	// handler timeout
	// handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	handlerTimeout := "5"
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	// time.Sleep(time.Second * 5)
	// ss.KillStreamerContainer(ctx)

	// Handler
	handler.NewHandler(&handler.Config{
		R:               router,
		InstanceService: instanceService,
		StreamerService: ss,
		BaseURL:         baseURL,
		TimeoutDuration: time.Duration(time.Duration(ht) * time.Second),
	})

	done := make(chan struct{})
	go func() {
		<-done
		fmt.Println("Closing containers.")

		// kill streamer
		ss.KillStreamerContainer(ctx)

		// kil instance containers
		instanceService.DeleteAllInstance(ctx)

		done <- struct{}{}
	}()

	return router, done, nil
}
