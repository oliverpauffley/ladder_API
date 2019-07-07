package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware(t *testing.T) {

	t.Run("Unauthorized users are rejected", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "testing", nil)
		response := httptest.NewRecorder()

		handler := AuthMiddleware(getVoidHandler())
		handler.ServeHTTP(response, req)

		want := http.StatusUnauthorized
		got := response.Code
		if want != got {
			t.Errorf("Expected %v got %v", want, got)
		}
	})

	t.Run("Authorized users are accepted", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/testing", nil)
		response := httptest.NewRecorder()

		// get valid token
		user := User{
			1,
			"ollie",
			jwt.StandardClaims{
				// token lasts 1 day
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
				Issuer:    "ladderapp",
			},
		}

		// create new token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

		tokenString, err := token.SignedString([]byte(SECRETKEY))
		if err != nil {
			t.Fatalf("Should be no error here, %v", err)
		}

		req.Header.Set("Authorization", tokenString)

		handler := AuthMiddleware(getTestHandler())
		handler.ServeHTTP(response, req)

		want := http.StatusOK
		got := response.Code
		if want != got {
			t.Errorf("Expected %v got %v", want, got)
		}
		wantString := "success!"
		gotString := response.Body.String()
		if gotString != wantString {
			t.Errorf("Inner handler does not run, wanted %s got %s", wantString, gotString)
		}
	})
}

// getVoid Handler returns a http.HandlerFunc for testing http middleware, should never run!
func getVoidHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		panic("test entered test handler, this should not happen")
	}
	return http.HandlerFunc(fn)
}

// GetTestHandler returns a http.HandlerFunc for testing http middleware
func getTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "success!")
	}
	return http.HandlerFunc(fn)
}
