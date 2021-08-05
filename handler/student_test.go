package handler

import (
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

var studentTestData = []models.StudentExam{
	{
		Exam:      1,
		StudentID: "test.person1",
		Score:     0.50,
	},
	{
		Exam:      2,
		StudentID: "test.person1",
		Score:     0.90,
	},
	{
		Exam:      1,
		StudentID: "test.person2",
		Score:     0.60,
	},
	{
		Exam:      1,
		StudentID: "test.person3",
		Score:     0.70,
	},
	{
		Exam:      2,
		StudentID: "test.person4",
		Score:     0.80,
	},
}

func addStudentTestRoutes() (*mux.Router, error) {
	err := db.InitDB(config.DBSchema)
	if err != nil {
		return nil, err
	}

	for _, exam := range studentTestData {
		err = db.UpsertRow(config.ScoreTable, exam)
		if err != nil {
			return nil, err
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/students", GetAllStudents).Methods("GET")
	router.HandleFunc("/students/{id}", GetStudentByID).Methods("GET")

	return router, nil
}

func TestGetAllStudents(t *testing.T) {
	router, err := addStudentTestRoutes()
	if err != nil {
		t.Errorf("Failed to start server")
	}

	request, _ := http.NewRequest("GET", "/students", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	want := 200
	if response.Code != want {
		t.Errorf("HTTP status is not OK; have: %v, want: %v", response.Code, want)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unable to read response body")
	}
	body := models.AllStudentListResponse{}
	err = json.Unmarshal(resBytes, &body)
	if err != nil {
		t.Errorf("Failed to parse response returned from route")
	}

	want = 4
	have := len(body.Students)
	if have != want {
		t.Errorf("Student count does not match expected value; have: %v, want: %v", have, want)
	}
}

func TestGetStudentByID(t *testing.T) {
	router, err := addStudentTestRoutes()
	if err != nil {
		t.Errorf("Failed to start server")
	}

	// Test with an invalid id
	request, _ := http.NewRequest("GET", "/students/does_not_exist", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Verify status code is 404
	have := response.Code
	want := 404
	if have != want {
		t.Errorf("Route returned an incorrect status code; have: %v, want: %v", have, want)
	}

	// Test a valid user
	request, _ = http.NewRequest("GET", "/students/test.person1", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// verify status code is 200
	have = response.Code
	want = 200
	if have != want {
		t.Errorf("Route returned an incorrect status code; have: %v, want: %v", have, want)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unable to read response body")
	}
	body := models.StudentByIDResponse{}
	err = json.Unmarshal(resBytes, &body)
	if err != nil {
		t.Errorf("Failed to parse response returned from route")
	}

	// Verify the correct number of exams are returned
	have = len(body.Exams)
	want = 2
	if have != want {
		t.Errorf("Incorrect number of exams returned; have: %v, want: %v", have, want)
	}

	// Verify the average is correct
	have64 := body.Average
	want64 := 0.7
	if have64 != want64 {
		t.Errorf("Incorrect average; have: %v, want: %v", have, want)
	}
}
