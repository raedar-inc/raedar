package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONSuccess Reserved field to add some meta information to the API response
type JSONSuccess struct {
	Meta    interface{} `json:"meta"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

// APIError shows the response structure for the programs api error
type APIError struct {
	Status  int16       `json:"status"`
	Error   interface{} `json:"error"`
	Success bool        `json:"success"`
}

// JSON returns a response for an API request
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// Message is a response returned to a user
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// ERROR response is used to return a JSON erro response to a client
func ERROR(w http.ResponseWriter, statusCode int, err *APIError) {
	// JSON(w, statusCode, err)
	if err != nil {
		JSON(w, statusCode, struct {
			Error interface{} `json:"error"`
		}{
			Error: err,
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)

	// w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(statusCode)
	// err := json.NewEncoder(w).Encode(err)
	// if err != nil {
	// 	fmt.Fprintf(w, "%s", err.Error())
	// }
}
