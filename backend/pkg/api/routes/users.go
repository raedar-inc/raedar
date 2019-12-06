package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"raedar/pkg/api/responses"
)

type userRegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register registers a new user into the system.
func Register() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		response := &responses.JSONSuccess{Data: "Welcome to raedar", Success: true}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}
}

// // UserRoutes creates all user routes
// func UserRoutes(router httprouter.Router) {

// 	router.GET("/api/v1/signup", Register())
// }

// type server struct{}

// func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ReponseWrite, r *http.HandlerFun) {
// 		if !currentUser(r).IsAdmin {
// 			http.NotFound(w, r)
// 			return
// 		}
// 		h(w, r)
// 	}
// }

// type server struct {
// 	db *someDatabase
// 	router *someRouter
// 	email EmailSender
// }

// func (*server) responde(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
// 	w.WriteHeader(status)
// 	if data != nil {
// 		err := json.NewEncoder(w).Encode(data)
// 	}
// }
