package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	// create new get request
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a test response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	// serve handler and req/reponse
	handler.ServeHTTP(rr, req)

	// check for correct status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returns the wrong status code, expected %v, got %v", http.StatusOK, rr.Code)
	}

	// check response body
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returns the wrong string, expected %v, got %v", expected, rr.Body.String())
	}
}
