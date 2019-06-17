package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux" // mux router for routes
	_ "github.com/lib/pq"    // postgresql driver
)

// database package level variable
var db *sql.DB

func main() {
	defer db.Close()
	r := NewRouter()
	// start database connection
	initDB()
	log.Fatal(http.ListenAndServe("localhost:8080", r))

}

// create a new router
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// set up routes
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/healthcheck", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	return r
}

// launch db
func initDB() {
	var err error
	// connect to postgres db
	db, err = sql.Open("postgres", "dbname=chess_ladder user=chess_admin sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}
