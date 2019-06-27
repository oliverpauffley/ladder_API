package main

import (
	"github.com/gorilla/sessions"
	"github.com/oliverpauffley/chess_ladder/models"
	"log"
	"net/http"

	_ "github.com/lib/pq" // postgresql driver
)

// env variable to store package environments variables
type Env struct {
	db models.Datastore
}

// gorilla sessions key set
var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func main() {
	// start database connection
	db, err := models.NewDB("postgres://chess_admin@localhost/chess_ladder?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db}
	router := env.NewRouter()
	log.Fatal(http.ListenAndServe("localhost:8080", router))

}
