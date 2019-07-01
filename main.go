package main

import (
	"encoding/gob"
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

// gorilla sessions cookie store
var store *sessions.CookieStore

// Pre-run setup
func init() {
	// register User struct with cookie store
	gob.Register(User{})

	// setup store with random key
	key := []byte("secret-key")
	store = sessions.NewCookieStore(key)

	// Set Cookies to last one day
	store.Options = &sessions.Options{MaxAge: 60 * 60 * 24}
}

func main() {
	// start database connection
	db, err := models.NewDB("postgres://chess_admin@localhost/chess_ladder?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db}
	Router := env.NewRouter()
	log.Fatal(http.ListenAndServe("localhost:8080", Router))

}
