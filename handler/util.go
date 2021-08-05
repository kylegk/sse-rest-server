package handler

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/kylegk/sse-rest-server/models"
)

// PanicRecovery is a middleware function used to inform the client an error has occurred and gracefully recover from a panic
func PanicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				fmt.Printf("recovering from error: %v\n %s", err, buf)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"Internal Server Error"}`))
			}
		}()

		h.ServeHTTP(w, r)
	})
}

// LogRequest is a middleware logger that will log basic information about every request
func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Verify that the request body (StudentExam struct) is valid
func validateRequestBody(exam models.StudentExam) error {
	if exam.Exam == 0 {
		return fmt.Errorf("invalid exam id")
	}
	if exam.Score == 0 {
		return fmt.Errorf("invalid score")
	}
	if exam.StudentID == "" {
		return fmt.Errorf("invalid studentid")
	}

	return nil
}
