package main

import "github.com/dgrijalva/jwt-go"

// the front end should send the following to login and register
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterCredentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

type AddLadderCredentials struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Owner  int    `json:"owner"`
}

type JoinLadderCredentials struct {
	HashId string `json:"hashid"`
	Id     int    `json:"id"`
}

// Custom jwt claims struct that inherits from standard
type User struct {
	ID       int    `json:"ID"`
	Username string `json:"username"`
	jwt.StandardClaims
}
