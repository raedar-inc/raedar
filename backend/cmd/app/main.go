package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// JSONResponse Reserved field to add some meta information to the API response
type JSONResponse struct {
	Meta    interface{} `json:"meta"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

// APIError shows the response structure for the programs api error
type APIError struct {
	Status  int16  `json:"status"`
	Code    int16  `json:"code"`
	Message string `json:"message"`
}

func main() {
	router := httprouter.New()

	// Route Handlers / Endpoints
	router.GET("/api/", indexHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	response := &JSONResponse{Data: "Welcome to raedar", Success: true}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// a small file glueing things together.
