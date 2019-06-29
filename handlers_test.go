package main

import (
	"bytes"
	"encoding/json"
	"github.com/oliverpauffley/chess_ladder/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegisterHandler(t *testing.T) {
	// create mock db and environment
	mdb := Mockdb{map[string]models.Credentials{}}
	env := Env{db: mdb}

	var tt = []struct {
		name  string
		input jsonCredentials
		want  int
	}{
		{"returns bad request when sent empty values",
			jsonCredentials{
				"",
				"",
				""},
			http.StatusBadRequest},

		{"returns bad request when passwords don't match",
			jsonCredentials{"pete", "hello", "goodbye"},
			http.StatusBadRequest},

		{"accepts a good json packet",
			jsonCredentials{"rob", "goodpassword", "goodpassword"},
			http.StatusOK},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			b, _ := json.Marshal(test.input)

			// form request and response
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(env.RegisterHandler)
			handler.ServeHTTP(response, req)

			got := response.Code
			if test.want != got {
				t.Errorf("Expected %v got %v", test.want, got)
			}
		})
	}

	t.Run("stops user registering when username already exists", func(t *testing.T) {

		input1 := jsonCredentials{"ollie", "1234", "1234"}
		input2 := jsonCredentials{"ollie", "hello", "hello"}

		a, _ := json.Marshal(input1)
		b, _ := json.Marshal(input2)

		// form two requests and response
		reqA, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(a))
		reqB, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(env.RegisterHandler)
		handler.ServeHTTP(response, reqA)

		want := http.StatusOK
		got := response.Code
		if want != got {
			t.Errorf("Expected %v got %v", want, got)
		}

		handler.ServeHTTP(response, reqB)

		want = http.StatusConflict
		got = response.Code
		if want != got {
			t.Errorf("Expected %v got %v", want, got)
		}
	})
}

func TestLoginHandler(t *testing.T) {
	// create mock db and environment
	users := make(map[string]models.Credentials)
	hash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	users["ollie"] = models.Credentials{Id: 1, Username: "ollie", JoinDate: time.Now(), Role: 1, Wins: 0, Losses: 0, Draws: 0, Hash: hash}
	mdb := Mockdb{users}
	env := Env{db: mdb}

	var tt = []struct {
		name  string
		input jsonCredentials
		want  int
	}{
		{"allows user to login",
			jsonCredentials{"ollie", "12345", ""},
			http.StatusOK},

		{"rejects users not in the db",
			jsonCredentials{"Paula", "12345", ""},
			http.StatusUnauthorized},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			b, _ := json.Marshal(test.input)

			// form request and response
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(env.LoginHandler)
			handler.ServeHTTP(response, req)

			got := response.Code
			if test.want != got {
				t.Errorf("Expected %v got %v", test.want, got)
			}
		})
	}
}
