package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

// Authenticate using gorilla sessions. wraps around http handlers
func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get jwt from request
		tokenString := r.Header.Get("Authorization")

		// If token is empty return unauthorized
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// get user instance to store payload
		user := &User{}

		// parse jwt string into user struct
		token, err := jwt.ParseWithClaims(tokenString, user, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtKey), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Print("signature invalid")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Printf("Error parsing signature, %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !token.Valid {
			log.Print("token not valid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// pass user into request
		context.Set(r, "decoded", user)
		// auth passed so continue to handler
		h.ServeHTTP(w, r)
	})
}
