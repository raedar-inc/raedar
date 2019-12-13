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
		if err != nil {
			h.logger.Print(err)
		}
		h.logger.Print(data)
		var errResponse = &responses.APIError{}
		errResponse = &responses.APIError{Error: err, Success: true, Status: http.StatusBadRequest}

		user := models.User{}
		userService := services.User{}
		err = json.Unmarshal(data, &user)
		if err != nil {
			errResponse = &responses.APIError{Error: err, Success: true, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			h.logger.Print(err)
			return
		}

		userCreated, erro := userService.Save(&user)
		if erro != "" {
			h.logger.Print(erro)
			errResponse = &responses.APIError{Error: erro, Success: false, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			return
		}

		response := &responses.JSONSuccess{Data: userCreated, Success: true}
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
}
