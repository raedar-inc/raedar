package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"raedar/pkg/api/responses"
	"raedar/pkg/repository/models"
	"raedar/pkg/repository/services"
)

type userRegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handlers for all rest handlers
type Handlers struct {
	logger *log.Logger
}

// Register registers a new user into the system.
func (h *Handlers) register() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := ioutil.ReadAll(r.Body)
		var errResponse = &responses.APIError{}
		if err != nil {
			h.logger.Print(err)
			errResponse = &responses.APIError{Error: err, Success: false, Status: http.StatusBadRequest}
		}

		user := models.User{}
		userService := services.User{}
		err = json.Unmarshal(data, &user)
		if err != nil {
			errResponse = &responses.APIError{Error: err, Success: false, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			h.logger.Print(err)
			return
		}

		userCreated, erro := userService.Save(&user)
		if erro != "" {
			errResponse = &responses.APIError{Error: erro, Success: false, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			return
		}

		response := &responses.JSONSuccess{Data: userCreated, Success: true}
		responses.JSON(w, http.StatusOK, response)
	}
}

// Login a user into the system.
func (h *Handlers) login() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := ioutil.ReadAll(r.Body)
		var errResponse = &responses.APIError{}
		var erro string

		if err != nil {
			h.logger.Print(err)
			errResponse = &responses.APIError{Error: err, Success: true, Status: http.StatusBadRequest}
		}

		user := &models.User{}
		userData := &models.User{}
		userService := services.User{}
		err = json.Unmarshal(data, &user)
		err = json.Unmarshal(data, &userData)
		if err != nil {
			errResponse = &responses.APIError{
				Error:   "please provide email and password to login",
				Success: true,
				Status:  http.StatusBadRequest,
			}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			h.logger.Print(err)
			return
		}

		erro = userService.Validate("login", user)
		if erro != "" {
			h.logger.Print(erro)
			errResponse = &responses.APIError{Error: erro, Success: false, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			return
		}

		user, err = userService.FindByEmail(user.Email)
		if err != nil {
			errResponse = &responses.APIError{Error: "No user found", Success: false, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			return
		}

		err = userService.VerifyPassword(user.Password, userData.Password)
		if err != nil {
			errResponse = &responses.APIError{Error: "Wrong password provided", Success: false, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			return
		}

		accessToken, err := userService.AccessToken(user)
		refreshToken, err := userService.RefreshToken(user)
		if err != nil {
			h.logger.Print(err)
			errResponse = &responses.APIError{Error: "Server error", Success: false, Status: http.StatusInternalServerError}
			responses.ERROR(w, http.StatusInternalServerError, errResponse)
			return
		}

		response := &responses.JSONSuccess{
			Data: map[string]interface{}{
				"email":         user.Email,
				"username":      user.Username,
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
			Success: true,
		}
		responses.JSON(w, http.StatusOK, response)
	}
}

// NewHandler returns user handlers struct
func NewHandler(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

// Routes sets up authentication Routes
func (h *Handlers) Routes(router *httprouter.Router) {
	router.POST("/api/v1/signup", h.register())
	router.POST("/api/v1/login", h.login())
}
