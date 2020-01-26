package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type response struct {
	Success bool  `json:"success"`
	ID      int64 `json:"id"`
}

func main() {
	randID := rand.New(rand.NewSource(time.Now().UnixNano()))

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = "0.0.0.0:8081"
	}
	r := mux.NewRouter()

	r.HandleFunc("/api/orders/create", func(w http.ResponseWriter, r *http.Request) {

		log.Println("Received payload")
		respondWithJSON(w, http.StatusOK, response{
			true,
			randID.Int63(),
		})
	}).Methods("POST")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGTERM)

	stopServerChan := make(chan struct{})

	go func() {

		srv := &http.Server{
			Handler: r,
			Addr:    httpAddr,
			// Set timeouts to avoid Slowloris attacks.
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())

	}()

	<-stop
	log.Println("shutting down gracefully")

	stopServerChan <- struct{}{}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
