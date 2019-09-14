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

// struct for environment variables
type Config struct {
	DbUser, DbPassword, JwtKey, HashKey string
}

var config Config

func main() {
	//  load env variables and make connection string
	config = getConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@db/chess_ladder?sslmode=disable",
		config.DbUser, config.DbPassword)

	// start database connection
	db, err := models.NewDB(connStr)
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db}
	Router := env.NewRouter()

	//use cors to manage cross origin requests
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8080"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	// set cors to handle all requests
	handler := c.Handler(Router)

	log.Fatal(http.ListenAndServe(":8000", handler))

}

func getConfig() Config {
	JwtKey, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		log.Panic("JWT_KEY unset, please set this key before running ")
	}
	HashKey, ok := os.LookupEnv("HASH_KEY")
	if !ok {
		log.Panic("HASHKEY unset, please set this key before running ")
	}
	dbPassword, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		log.Panic("POSTGRES_PASSWORD unset, please set this key before running ")
	}
	dbUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		log.Panic("POSTGRES_USER unset, please set this key before running ")
	}
	return Config{
		DbUser:     dbUser,
		DbPassword: dbPassword,
		JwtKey:     JwtKey,
		HashKey:    HashKey,
	}
}
