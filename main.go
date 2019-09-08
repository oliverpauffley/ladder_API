package main

import (
	"fmt"
	_ "github.com/lib/pq" // postgresql driver
	"github.com/oliverpauffley/chess_ladder/models"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

// env variable to store package environments variables
type Env struct {
	db models.Datastore
}

//JWT secret key, change on prod!
const SECRETKEY string = "secret"

func main() {
	//  load env variables and make connection string
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	connStr := fmt.Sprintf("postgres://%s:%s@db/chess_ladder?sslmode=disable",
		dbUser, dbPassword)

	// start database connection
	db, err := models.NewDB(connStr)
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db}
	log.Printf("Connected to db")
	Router := env.NewRouter()

	//use cors to manage cross origin requests
	// change these options on prod
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8080"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	// set cors to handle all requests
	handler := c.Handler(Router)

	log.Fatal(http.ListenAndServe("localhost:8000", handler))

}
