package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// create a new router
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// set up routes
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/healthcheck", HealthCheckHandler).Methods("GET")

	return r
}

func main() {
	r := NewRouter()

	log.Fatal(http.ListenAndServe("localhost:8080", r))

}
