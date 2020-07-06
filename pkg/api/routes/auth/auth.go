package auth

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"raedar/pkg/utils"

	"github.com/julienschmidt/httprouter"

	"raedar/pkg/api/responses"
	"raedar/pkg/repository/models"
	"raedar/pkg/repository/services"
)

// Handlers for all rest handlers
type Handlers struct {
	logger *log.Logger
}

var (
	userService = services.User{}
	errResponse = &responses.APIError{}
)

func passwordHandler(h *Handlers, r *http.Request, params httprouter.Params) (*responses.APIError, int) {
	token := params.ByName("token")
	data, err := utils.DecodeToken(token)

	type requestData struct {
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err != nil {
		h.logger.Print(err)
		return &responses.APIError{
			Error:   "Your token has expired",
			Success: false,
			Status:  http.StatusBadRequest,
		}, http.StatusUnprocessableEntity
	}
	email, ok := data.(jwt.MapClaims)["email"]
	emailStr, ok := email.(string)
	if !ok {
		return &responses.APIError{
			Error:   "Invalid token",
			Success: false,
			Status:  http.StatusBadRequest,
		}, http.StatusUnprocessableEntity
	}

	userData := requestData{}
	bodyData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(bodyData, &userData)
	if err != nil {
		return &responses.APIError{
			Error:   "please an password and a confirm password",
			Status:  http.StatusBadRequest,
			Success: false,
		}, http.StatusUnprocessableEntity
	}
	user, err := userService.FindByEmail(emailStr)
	if err == nil {
		err = userService.VerifyPassword(user.Password, userData.Password)
		if err == nil {
			h.logger.Print(err)
			return &responses.APIError{
				Error:   "password cannot be the same as current password",
				Status:  http.StatusBadRequest,
				Success: false,
			}, http.StatusUnprocessableEntity
		}
		if !userService.ComparePasswordToConfirmPassword(userData.Password, userData.ConfirmPassword) {
			return &responses.APIError{
				Error:   "password and confirmPassword do not match",
				Status:  http.StatusBadRequest,
				Success: false,
			}, http.StatusBadRequest
		}
		hashedPassword, err := services.HashPassword(userData.Password)
		if err != nil {
			return &responses.APIError{
				Error:   "Something went wrong",
				Status:  http.StatusUnprocessableEntity,
				Success: false,
			}, http.StatusUnprocessableEntity
		}
		user.Password = string(hashedPassword)
		_, err = userService.Update(user, map[string]interface{}{"password": string(hashedPassword)})
		if err != nil {
			return &responses.APIError{
				Error:   "something went wrong",
				Status:  http.StatusUnprocessableEntity,
				Success: false,
			}, http.StatusUnprocessableEntity
		}
	}
	return nil, 0
}

// Register registers a new user into the system.
func (h *Handlers) register() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.logger.Print(err)
			errResponse = &responses.APIError{Error: err, Success: false, Status: http.StatusBadRequest}
		}

		user := models.User{}
		err = json.Unmarshal(data, &user)
		if err != nil {
			errResponse = &responses.APIError{Error: err, Success: false, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			h.logger.Print(err)
			return
		}

		userCreated, errStr := userService.Save(&user)
		if errStr != "" {
			errResponse = &responses.APIError{Error: errStr, Success: false, Status: http.StatusBadRequest}
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
		var errStr string

		if err != nil {
			errResponse = &responses.APIError{Error: err, Success: false, Status: http.StatusBadRequest}
		}

		user := &models.User{}
		userData := &models.User{}
		err = json.Unmarshal(data, &user)
		err = json.Unmarshal(data, &userData)
		if err != nil {
			errResponse = &responses.APIError{
				Error:   "please provide email and password to login",
				Success: false,
				Status:  http.StatusBadRequest,
			}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			h.logger.Print(err)
			return
		}

		errStr = userService.Validate("login", user)
		if errStr != "" {
			errResponse = &responses.APIError{Error: errStr, Success: false, Status: http.StatusBadRequest}
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

// Request forgot password functionality for users who forgot their passwords.
func (h *Handlers) requestForgotPassword() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var errResponse = &responses.APIError{}
		var errStr string
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.logger.Print(err)
			errResponse = &responses.APIError{Error: err, Success: false, Status: http.StatusBadRequest}
		}

		userData := &models.User{}
		userService := services.User{}
		err = json.Unmarshal(data, &userData)
		if err != nil {
			errResponse = &responses.APIError{
				Error:   "please provide an email",
				Status:  http.StatusBadRequest,
				Success: false,
			}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			h.logger.Print(err)
			return
		}

		errStr = userService.Validate("reset-password", userData)
		if errStr != "" {
			errResponse = &responses.APIError{Error: errStr, Success: false, Status: http.StatusBadRequest}
			responses.ERROR(w, http.StatusUnprocessableEntity, errResponse)
			return
		}

		if user, _ := userService.FindByEmail(userData.Email); user != nil {
			var resetPasswordUrl bytes.Buffer
			var msg bytes.Buffer
			resetPasswordUrl.WriteString("http://127.0.0.1:8080/api/v1/password/forgot/")
			token, _ := utils.CreateResetPasswordToken(userData.Email)
			resetPasswordUrl.WriteString(token)
			msg.WriteString(
				`Click the link below to reset your password, if this wasn't you, please ignore this message`)
			msg.WriteString("\n\n")
			msg.WriteString(resetPasswordUrl.String())
			msg.WriteString("\n\n")
			msg.WriteString("Thanks")
			emailSubject := "Reset Forgotten Password"
			email := utils.Email{Email: userData.Email}
			emailErr := make(chan error)
			go email.SendEmail(emailSubject, msg.String(), emailErr)
			h.logger.Print(<-emailErr)
		}

		response := &responses.JSONSuccess{
			Data: map[string]string{
				"message": "a password-reset url has been sent to your email to reset your password",
			},
			Success: true,
		}
		responses.JSON(w, http.StatusOK, response)
	}
}

func (h *Handlers) forgotPassword() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		errResponse, status := passwordHandler(h, r, params)
		if errResponse != nil {
			responses.ERROR(w, status, errResponse)
			return
		}

		response := &responses.JSONSuccess{
			Data: map[string]string{
				"message": "Your password has been reset successfully, you can now login with your new password",
			},
			Success: true,
		}
		responses.JSON(w, http.StatusOK, response)
	}
}

// Reset password functionality for users who forgot their passwords.
func (h *Handlers) resetPassword() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		errResponse, status := passwordHandler(h, r, params)
		if errResponse != nil {
			responses.ERROR(w, status, errResponse)
			return
		}

		token := params.ByName("token")
		data, _ := utils.DecodeToken(token)
		email, _ := data.(jwt.MapClaims)["email"]

		response := &responses.JSONSuccess{
			Data: map[string]interface{}{
				"email":   email,
				"message": "Your password has been reset successfully",
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
	router.POST("/api/v1/password/reset/:token", h.resetPassword())
	router.POST("/api/v1/password/forgot/:token", h.forgotPassword())
	router.POST("/api/v1/password/forgot-request", h.requestForgotPassword())
}
