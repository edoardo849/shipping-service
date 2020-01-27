package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/edoardo849/bezos/pkg/api"
	"github.com/edoardo849/bezos/pkg/delivery"
	"github.com/edoardo849/bezos/pkg/order"
	"github.com/gorilla/mux"
)

func main() {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = "0.0.0.0:8080"
	}

	dbConn := os.Getenv("DB_CONN")
	if dbConn == "" {
		dbConn = "docker:docker@tcp(db:3306)/api_db"
	}

	deliveryAddr := os.Getenv("DELIVERY_ADDR")
	if deliveryAddr == "" {
		deliveryAddr = "docker:docker@tcp(db:3306)/api_db"
	}

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGTERM)

	stopServerChan := make(chan struct{})

	orderChan := make(chan order.CreateReq, 1)
	orderService := order.New(db, orderChan)
	delivery := delivery.New(
		&http.Client{
			Timeout: 5 * time.Second,
		},
		deliveryAddr,
		orderService,
	)

	stopDeliveryService := delivery.Start(orderChan)

	server := api.New(
		orderService,
		mux.NewRouter(),
		stopServerChan,
	)

	go func() {
		srv := &http.Server{
			Addr: httpAddr,
			// Set timeouts to avoid Slowloris attacks.
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		}

		if err := server.ServeHTTP(srv); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-stop
	log.Println("shutting down gracefully")
	stopDeliveryService()
	stopServerChan <- struct{}{}
}
