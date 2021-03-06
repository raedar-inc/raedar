package routes

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"context"

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
	s := http.Server{
		Addr:              ":8080",
		Handler:           router,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
	}

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	<-ch

	ctx, ctxFunc := context.WithTimeout(context.Background(), 100*time.Second)
	ctxFunc()
	s.Shutdown(ctx)
}
