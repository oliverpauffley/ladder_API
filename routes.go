package main

import "github.com/gorilla/mux"

// create a new router
func (env *Env) NewRouter() *mux.Router {
	router := mux.NewRouter()

	// set up routes
	router.HandleFunc("/register", env.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", env.LoginHandler).Methods("POST")

	// create authenticated routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Use(AuthMiddleware)
	authRouter.HandleFunc("/logout", env.LogoutHandler).Methods("GET")

	// user routes
	authRouter.HandleFunc("/password", env.ChangePasswordHandler).Methods("POST")

	// ladder routes
	authRouter.HandleFunc("/users/{id:[0-9]+}", env.UserHandler).Methods("GET")
	authRouter.HandleFunc("/ladder", env.AddLadderHandler).Methods("POST")
	authRouter.HandleFunc("/ladder/user/{id:[0-9]+}", env.GetAllLaddersHandler).Methods("GET")
	authRouter.HandleFunc("/ladder/join", env.JoinLadderHandler).Methods("POST")

	// game routes
	authRouter.HandleFunc("/game", env.AddGame).Methods("POST")

	return router
}
