package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/oliverpauffley/chess_ladder/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestRegisterHandler(t *testing.T) {
	// create mock db and environment
	mdb := Mockdb{map[string]models.CredentialsInternal{}}
	env := Env{db: mdb}

	var tt = []struct {
		name  string
		input RegisterCredentials
		want  int
	}{
		{"returns bad request when sent empty values",
			RegisterCredentials{
				"",
				"",
				"",
				""},
			http.StatusBadRequest,
		},

		{"returns bad request when passwords don't match",
			RegisterCredentials{"pete", "pete@example.com", "goodbye", "hello"},
			http.StatusBadRequest,
		},

		{"accepts a good json packet",
			RegisterCredentials{"rob", "rob@example.com", "goodpassword", "goodpassword"},
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

		input1 := RegisterCredentials{"ollie", "ollie@example.com", "1234", "1234"}
		input2 := RegisterCredentials{"ollie", "ollie@example.com", "hello", "hello"}

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
	users := make(map[string]models.CredentialsInternal)
	hash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	users["ollie"] = models.CredentialsInternal{Id: 1, Username: "ollie", JoinDate: time.Now(), Role: 1, Wins: 0, Losses: 0, Draws: 0, Hash: hash}
	mdb := Mockdb{users}
	env := Env{db: mdb}

	var tt = []struct {
		name  string
		input LoginCredentials
		want  int
	}{
		{"allows user to login",
			LoginCredentials{"ollie", "12345"},
			http.StatusOK},

		{"rejects users not in the db",
			LoginCredentials{"Paula", "12345"},
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

func TestUserHandler(t *testing.T) {
	// create mock db and environment
	users := make(map[string]models.CredentialsInternal)
	hash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	users["ollie"] = models.CredentialsInternal{Id: 1, Username: "ollie", JoinDate: time.Now(), Role: 1, Wins: 0, Losses: 0, Draws: 0, Hash: hash}
	mdb := Mockdb{users}
	env := Env{db: mdb}

	var tt = []struct {
		name        string
		ID          int
		want        models.CredentialsExternal
		code        int
		contentType string
	}{
		{"Shows stats for a user in db",
			1,
			models.CredentialsExternal{Id: users["ollie"].Id, Username: users["ollie"].Username,
				JoinDate: users["ollie"].JoinDate.Round(time.Hour), Role: users["ollie"].Role, Wins: users["ollie"].Wins,
				Losses: users["ollie"].Losses, Draws: users["ollie"].Draws},
			http.StatusOK,
			"application/json",
		},

		{"sends empty user for user not in db with error",
			4,
			models.CredentialsExternal{},
			http.StatusNotFound,
			"application/json",
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			// create request and response with correct url variables
			urlString := fmt.Sprintf("auth/users/%d", test.ID)
			req, _ := http.NewRequest(http.MethodGet, urlString, nil)
			response := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(test.ID)})

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

			// write token to header
			req.Header.Set("Authorization", tokenString)

			handler := AuthMiddleware(http.HandlerFunc(env.UserHandler))
			handler.ServeHTTP(response, req)

			// decode and store request json
			got := models.CredentialsExternal{}
			if response.Code != http.StatusNotFound {
				err = json.Unmarshal(response.Body.Bytes(), &got)
				if err != nil {
					t.Fatalf("Error decoding json, err %v, json body %v", err.Error(), response.Body.Bytes())
				}
			}

			// compare the response using go-cmp package as reflect.deepequal fails
			if !cmp.Equal(got, test.want) {
				t.Errorf("Did not get back correct user credentials, got %v, want %v", got, test.want)
			}

			if response.Code != test.code {
				t.Errorf("Got the wrong response code, got %d, wanted %d", response.Code, test.code)
			}
		})
	}
}
