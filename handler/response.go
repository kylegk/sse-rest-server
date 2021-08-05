package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kylegk/sse-rest-server/models"
)

// SendGenericNotFoundResponse returns a generic 404 error
func SendGenericNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(&models.GenericResponse{Error: "Not Found", Code: http.StatusNotFound, Message: "Resource not found"}, http.StatusNotFound, w)
}

// SendGenericNotAllowedResponse returns a generic 405 error
func SendGenericNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(&models.GenericResponse{Error: "Not Allowed", Code: http.StatusMethodNotAllowed, Message: "You are not authorized to access this resource"}, http.StatusMethodNotAllowed, w)
}

// SendGenericInternalServerError returns a generic 500 error
func SendGenericInternalServerError(w http.ResponseWriter, r *http.Request) {
	sendResponse(&models.GenericResponse{Error: "Internal Server Error", Code: http.StatusInternalServerError, Message: "An error has occurred"}, http.StatusInternalServerError, w)
}

// sendResponse is a generic method to send a custom response to the client
func sendResponse(payload interface{}, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
