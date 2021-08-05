package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kylegk/sse-rest-server/config"
	"github.com/kylegk/sse-rest-server/db"
	"github.com/kylegk/sse-rest-server/handler"
	"github.com/kylegk/sse-rest-server/sse"
	"log"
	"net/http"
	"os"
)

// Init creates the database, adds the routes to be handled and performs any other initial setup required
func Init(c config.Config) {
	var err error
	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()

	err = validateConfig(c)
	if err != nil {
		log.Println(err)
		return
	}

	err = db.InitDB(c.MemDBSchema)
	if err != nil {
		log.Println(err)
		return
	}

	sse.IngestData(c.SSEServerUrl)
	addRoutes(c.PORT)
}

func validateConfig(c config.Config) error {
	if c.PORT == "" || c.SSEServerUrl == "" || c.MemDBSchema == nil {
		return fmt.Errorf("invalid configuration")
	}

	return nil
}

// Initialize the routes to be served and setup any middleware applied to the routes
func addRoutes(port string) {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(handler.SendGenericNotFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(handler.SendGenericNotAllowedResponse)

	// Student route handlers
	router.HandleFunc("/students", handler.GetAllStudents).Methods("GET")
	router.HandleFunc("/students/{id}", handler.GetStudentByID).Methods("GET")

	// Exam route handlers
	router.HandleFunc("/exams", handler.GetAllUniqueExamIDs).Methods("GET")
	router.HandleFunc("/exams/all", handler.GetAllExams).Methods("GET")
	router.HandleFunc("/exams/{id}", handler.GetExamByID).Methods("GET")
	router.HandleFunc("/exams/{id}", handler.DeleteExam).Methods("DELETE")
	router.HandleFunc("/exams", handler.AddExam).Methods("POST")

	// Add panic middleware
	router.Use(handler.PanicRecovery)

	log.Fatal(http.ListenAndServe(port, handler.LogRequest(router)))
}
