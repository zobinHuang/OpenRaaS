package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info(os.Getenv("SERVICE_NAME"), " start to run")

	// injection
	router, err := inject()
	if err != nil {
		log.Fatalf("Failed to conduct injection:", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// provider worker node online
	addr := "http://" + os.Getenv("SCHEDULER_WS_HOSTNAME") + ":" + os.Getenv("SCHEDULER_WS_PORT") + "/api/scheduler/node_online" + "?type=provider"
	log.Println("The provider worker node's info is sent to the scheduler's HTTP interface: " + addr)
	client := &http.Client{}
	data := make(map[string]interface{})
	data["id"] = os.Getenv("SERVER_ID")
	data["ip"] = os.Getenv("SERVER_IP")
	data["is_contain_gpu"], _ = strconv.ParseBool(os.Getenv("SERVER_PERFORMANCE"))
	data["processor"], _ = strconv.ParseFloat(os.Getenv("SERVER_PROCESSOR"), 64)
	data["bandwidth"], _ = strconv.ParseFloat(os.Getenv("SERVER_BW"), 64)
	data["latency"], _ = strconv.ParseFloat(os.Getenv("SERVER_LATENCY"), 64)
	// history
	// data["inst_history"] = os.Getenv("SERVER_HISTORY")

	bytesData, _ := json.Marshal(data)
	trimmedBytesData := bytesData[:len(bytesData)-1]
	hisData, _ := json.Marshal(os.Getenv("SERVER_HISTORY"))
	trimmedHisData := hisData[1 : len(hisData)-1]
	trimmedHisData = bytes.ReplaceAll(trimmedHisData, []byte("\\"), []byte(""))
	bytesData = []byte(string(trimmedBytesData) + ",\"inst_history\":" + string(trimmedHisData) + "}")

	req, _ := http.NewRequest("POST", addr, bytes.NewReader(bytesData))
	resp, _ := client.Do(req)
	log.Println("Sent provider worker node online with info:", string(bytesData))
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Get answer from the scheduler:", string(body))

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

	// Shutdown server
	log.Info("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: ", err)
	}
}
