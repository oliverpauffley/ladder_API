package main

import (
	"bytes"
	"encoding/json"
	"github.com/oliverpauffley/chess_ladder/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockdb struct {
	users map[string]models.Credentials
}

func (db mockdb) CreateUser(username, password string) error {
	var mockcredentials models.Credentials
	mockcredentials.Username = username
	mockcredentials.Hash = password
	db.users[username] = mockcredentials
	return nil
}

func (db mockdb) QueryByName(username string) (models.Credentials, error) {
	mockcredentials := db.users[username]
	return mockcredentials, nil
}

func TestRegisterHandler(t *testing.T) {
	// create mock db and environment
	mdb := mockdb{map[string]models.Credentials{}}
	env := Env{db: mdb}

	t.Run("returns bad request when empty json", func(t *testing.T) {
		// create bad request json
		jsonstr := []byte(``)

		// form request and response
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonstr))
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(env.RegisterHandler)
		handler.ServeHTTP(response, req)

		want := http.StatusBadRequest
		got := response.Code
		if want != got {
			t.Errorf("Expected %v got %v", want, got)
		}
	})

	t.Run("returns bad request when not sent correct json", func(t *testing.T) {
		// create bad request json

		jsonstr := jsonCredentials{
			Username: "",
			Password: "",
			Confirm:  "",
		}
		b, _ := json.Marshal(jsonstr)

		// form request and response
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(env.RegisterHandler)
		handler.ServeHTTP(response, req)

		want := http.StatusBadRequest
		got := response.Code
		if want != got {
			t.Errorf("Expected %v got %v", want, got)
		}
	})
}
