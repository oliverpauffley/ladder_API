package main

import "net/http"

// Authenticate using gorilla sessions. wraps around http handlers
func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get cookie session from store
		session, err := store.Get(r, "authentication-cookie")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// check for using authentication
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// user is authenticated so return original handler
		h.ServeHTTP(w, r)
	})
}
