package main

import (
	"business/dal"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info(os.Getenv("SERVICE_NAME"), " start to run")

	// initialize data sources
	ds, err := dal.InitDS()
	if err != nil {
		log.Fatalf("Failed to initialize data sources:", err)
	}

	// injection
	router, err := inject(ds)
	if err != nil {
		log.Fatalf("Failed to conduct injection:", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize web server: ", err)
		}
	}()
	log.WithFields(log.Fields{
		"Port": srv.Addr,
	}).Info("Web server start listening")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources
	if err := dal.CloseDS(ds); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: ", err)
	}

	// Shutdown server
	log.Info("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: ", err)
	}
}
