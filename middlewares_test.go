package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {

	t.Run("Unauthorized users are rejected", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "testing", nil)
		response := httptest.NewRecorder()

		handler := AuthMiddleware(getVoidHandler())
		handler.ServeHTTP(response, req)

		want := http.StatusForbidden
		got := response.Code
		if want != got {
			t.Errorf("Expected %v got %v", want, got)
		}
	})

	t.Run("Authorized users are accepted", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/testing", nil)
		response := httptest.NewRecorder()

		// set user as authenticated
		session, err := store.Get(req, "authentication-cookie")
		if err != nil {
			t.Fatal("should be no error here")
		}
		user := User{ID: 1, Authenticated: true}
		session.Values["user"] = user
		err = session.Save(req, response)
		if err != nil {
			t.Fatal("should be no error here")
		}

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
