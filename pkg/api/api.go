package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/edoardo849/bezos/pkg/order"

	"github.com/gorilla/mux"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// New Creates a new handler
func New(os order.Service, r *mux.Router, stopChan chan struct{}) *Server {
	return &Server{
		orderService: os,
		router:       mux.NewRouter(),
		stopChan:     stopChan,
	}
}

// Server is the server
type Server struct {
	orderService order.Service
	router       *mux.Router
	stopChan     chan struct{}
	http         *http.Server
}

// Run runs the server
func (s *Server) ServeHTTP(http *http.Server) error {
	// Initialize routes
	s.http = http

	// http://zabana.me/notes/enable-cors-in-go-api.html
	s.registerHandlers()
	s.http.Handler = s.router
	log.Printf("Server listening on %s\n", s.http.Addr)

	go func() {
		<-s.stopChan
		log.Println("Shutting down server")
		s.http.Shutdown(context.Background())
	}()

	return s.http.ListenAndServe()
}

// Register routes
func (s *Server) registerHandlers() {

	// Use gorilla/mux for rich routing.
	// See http://www.gorillatoolkit.org/pkg/mux
	r := s.router.PathPrefix(fmt.Sprintf("/%s", apiVersion)).Subrouter()

	r.HandleFunc("/orders", withBasicAuth(handleOrdersCreate(s.orderService))).Methods("POST")

	r.NotFoundHandler = handle404()

}
