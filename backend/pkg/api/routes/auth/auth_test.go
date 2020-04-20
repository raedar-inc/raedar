package auth

import (
	"bytes"
	"io/ioutil"
	"log"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "raedar/pkg/repository/engines"

	"github.com/julienschmidt/httprouter"
)

func TestLoginUser(t *testing.T) {
	router := httprouter.New()
	_, err := http.NewRequest("POST", "localhost:8080/api/v1/login", nil)
	if err != nil {
		t.Fatalf("Could not make a request: %v", err)
	}

	logger := log.New(os.Stdout, "raedar", log.LstdFlags|log.Lshortfile)
	authentication := NewHandler(logger)
	authentication.Routes(router)

	// a recorder
	rec := httptest.NewRecorder()
	res := rec.Result()
	if res.StatusCode == http.StatusOK {
		t.Errorf("Expected status OK but got: %v", res.StatusCode)
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
}

func TestRouting(t *testing.T) {
	router := httprouter.New()
	_, err := http.NewRequest("POST", "localhost:8080/api/v1/login", nil)
	if err != nil {
		t.Fatalf("Could not make a request: %v", err)
	}

	logger := log.New(os.Stdout, "raedar", log.LstdFlags|log.Lshortfile)
	authentication := NewHandler(logger)
	authentication.Routes(router)

	res, err := http.Post("localhost:8080/api/v1/login", "application/json", nil)
	if err != nil {
		t.Fatalf("Could not make a request: %v", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read body: %v", err)
	}

	if msg := string(bytes.TrimSpace(body)); msg != "Go home" {
		t.Errorf("Test failed because =: %v", err)
	}
	if res.StatusCode == http.StatusOK {
		t.Errorf("Expected status OK but got: %v", res.StatusCode)
	}
}
