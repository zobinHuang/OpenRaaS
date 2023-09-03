package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"serverd/dal"
	"syscall"
	"time"
)

func main() {
	fmt.Printf("Start!")

	// initialize data sources
	ds, err := dal.InitDS()
	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	router, done, err := inject(ds)
	if err != nil {
		log.Fatalf("Failure to conduct injection: %v\n", err)
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:3080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			done <- struct{}{}
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()
	log.Printf("Listening on port %v\n", srv.Addr)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// shutdown all built containers
	done <- struct{}{}
	// time.Sleep(time.Millisecond)
	<-done
	fmt.Println("Stoped all containers.")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	timectx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources
	if err := dal.CloseDS(ds); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(timectx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
