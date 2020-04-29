package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"raedar/pkg/api/routes/auth"
)

// Endpoints function handling of all routes/endpoints
func Endpoints() {
	router := httprouter.New()

	// Setting global options here to support CORS.
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", r.Header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}
		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	logger := log.New(os.Stdout, "raedar", log.LstdFlags|log.Lshortfile)
	authentication := auth.NewHandler(logger)

	// Route Handlers / Endpoints
	authentication.Routes(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
