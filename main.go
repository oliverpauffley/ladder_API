package main

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/oliverpauffley/chess_ladder/models"
	"github.com/rs/cors"
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

	//use cors to manage cross origin requests
	// change these options on prod
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost*"},
		AllowCredentials: true,
	})

	handler := c.Handler(Router)

	log.Fatal(http.ListenAndServe("localhost:8000", handler))

}
