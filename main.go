package main

import (
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

//JWT secret key, change on prod!
const SECRETKEY string = "secret"

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
		AllowedOrigins:   []string{"http://127.0.0.1:8080"},
		AllowCredentials: true,
	})
	// set cors to handle all requests
	handler := c.Handler(Router)

	log.Fatal(http.ListenAndServe("localhost:8000", handler))

}
