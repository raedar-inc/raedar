package auth

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	_ "raedar/pkg/repository/engines"

	"github.com/julienschmidt/httprouter"
)

func TestRouting(t *testing.T) {
	var err error
	router := httprouter.New()
	jsonStr := []byte(`{"username":"google", "password":"krs1krs1", "email":"gool@dol.io"}`)

	data := url.Values{}
	data.Set("email", "bar@lop.co")
	data.Set("username", "foo")
	data.Set("password", "bar")

	logger := log.New(os.Stdout, "raedar", log.LstdFlags|log.Lshortfile)
	authentication := NewHandler(logger)
	authentication.Routes(router)

	req, err := http.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	defer req.Body.Close()

	if err != nil {
		t.Fatalf("Could not make a request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v but got: %v", http.StatusOK, rr.Code)
	}
}
