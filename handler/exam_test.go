package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kylegk/sse-rest-server/config"
	"github.com/kylegk/sse-rest-server/db"
	"github.com/kylegk/sse-rest-server/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var examTestData = []models.StudentExam{
	{
		Exam:      1,
		StudentID: "test.person",
		Score:     0.67,
	},
	{
		Exam:      1,
		StudentID: "test.person2",
		Score:     0.75,
	},
	{
		Exam:      1,
		StudentID: "test.person3",
		Score:     0.98,
	},
	{
		Exam:      2,
		StudentID: "test.person",
		Score:     0.89,
	},
}

func addExamTestRoutes() (*mux.Router, error) {
	err := db.InitDB(config.DBSchema)
	if err != nil {
		return nil, err
	}

	for _, exam := range examTestData {
		err = db.UpsertRow(config.ScoreTable, exam)
		if err != nil {
			return nil, err
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/exams", GetAllUniqueExamIDs).Methods("GET")
	router.HandleFunc("/exams/all", GetAllExams).Methods("GET")
	router.HandleFunc("/exams/{id}", GetExamByID).Methods("GET")
	router.HandleFunc("/exams/{id}", DeleteExam).Methods("DELETE")
	router.HandleFunc("/exams", AddExam).Methods("POST")

	return router, nil
}

func TestGetAllUniqueExamIDs(t *testing.T) {
	router, err := addExamTestRoutes()
	if err != nil {
		t.Errorf("Failed to start server")
	}

	request, _ := http.NewRequest("GET", "/exams", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Verify status code is 200
	have := response.Code
	want := 200
	if have != want {
		t.Errorf("HTTP status is not OK; have %v, want %v", response.Code, want)
	}

	body := models.AllUniqueExamsListResponse{}
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unable to read buffer")
	}
	err = json.Unmarshal(resBytes, &body)
	if err != nil {
		t.Errorf("Failed to parse response returned from route")
	}

	// Verify count of exams is correct
	have = len(body.Exams)
	want = 2
	if have != want {
		t.Errorf("Exam count does not match expected value; have: %v, want: %v", have, want)
	}
}

func TestGetAllExams(t *testing.T) {
	router, err := addExamTestRoutes()
	if err != nil {
		t.Errorf("Failed to start server")
	}

	request, _ := http.NewRequest("GET", "/exams/all", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Verify status code is 200
	have := response.Code
	want := 200
	if have != want {
		t.Errorf("HTTP status is not OK; have %v, want %v", response.Code, want)
	}

	body := models.AllExamsListResponse{}
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error reading response body")
	}
	err = json.Unmarshal(resBytes, &body)
	if err != nil {
		t.Errorf("Error parsing response body")
	}

	// Verify correct count of exams
	have = len(body.Exams)
	want = 4
	if have != want {
		t.Errorf("Exam count does not match expected value; have: %v, want: %v", have, want)
	}
}

func TestGetExamByID(t *testing.T) {
	router, err := addExamTestRoutes()
	if err != nil {
		t.Errorf("Failed to start server")
	}

	request, _ := http.NewRequest("GET", "/exams/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Verify status code is 200
	have := response.Code
	want := 200
	if have != want {
		t.Errorf("HTTP status is not OK; have %v, want %v", response.Code, want)
	}

	body := models.ExamByIDResponse{}
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error reading response body")
	}
	err = json.Unmarshal(resBytes, &body)
	if err != nil {
		t.Errorf("Error parsing response body")
	}

	// Verify correct number of exams was returned
	have = body.Exam
	want = 1
	if have != want {
		t.Errorf("Method returned the wrong exam; have %v, want %v", have, want)
	}

	// Verify correct number of scores was returned
	have = len(body.Scores)
	want = 3
	if have != want {
		t.Errorf("Mismatched number of scores returned; have %v, want %v", have, want)
	}

	// Verify that the average is correct
	have64 := body.Average
	want64 := 0.7999999999999999
	if have != want {
		t.Errorf("Incorrect average score; have '%v', want '%v", have64, want64)
	}
}

func TestAddExam(t *testing.T) {
	router, err := addExamTestRoutes()
	if err != nil {
		t.Errorf("Failed to start server")
	}

	exam := models.StudentExam{
		Exam:      1,
		StudentID: "test.person5",
		Score:     100,
	}

	j, err := json.Marshal(exam)
	if err != nil {
		t.Errorf("Cannot marshal request struct")
		return
	}

	request, _ := http.NewRequest("POST", "/exams", bytes.NewBuffer(j))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Verify status code is 200
	have := response.Code
	want := 200
	if have != want {
		t.Errorf("HTTP status is not OK; have %v, want %v", response.Code, want)
	}
}

func TestDeleteExam(t *testing.T) {
	router, err := addExamTestRoutes()
	if err != nil {
		t.Errorf("Failed to start server")
	}

	request, _ := http.NewRequest("DELETE", "/exams/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Verify status code is 200
	have := response.Code
	want := 200
	if have != want {
		t.Errorf("HTTP status is not OK; have %v, want %v", response.Code, want)
	}
}
