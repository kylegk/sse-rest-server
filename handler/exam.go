package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kylegk/sse-rest-server/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kylegk/sse-rest-server/db"
	"github.com/kylegk/sse-rest-server/models"
)

// AddExam adds a single exam (PUT) to the datastore
func AddExam(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			if err.Error() == "internal_server_error" {
				SendGenericInternalServerError(w, r)
			} else {
				sendResponse(&models.GenericResponse{Code: http.StatusBadRequest, Error: "Bad Request", Message: err.Error()}, http.StatusBadRequest, w)
			}

			return
		}
	}()

	exam := models.StudentExam{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(bytes, &exam)
	if err != nil {
		err = errors.New("unable to parse request")
		return
	}

	err = validateRequestBody(exam)
	if err != nil {
		return
	}

	err = db.UpsertRow(config.ScoreTable, exam)
	if err != nil {
		err = errors.New("internal_server_error")
		return
	}

	sendResponse(&models.GenericResponse{Message: fmt.Sprintf("Succesfully added exam: %v", exam.Exam)}, http.StatusOK, w)
}

// DeleteExam removes all examTestData matching the specified exam id
func DeleteExam(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			SendGenericInternalServerError(w, r)
			return
		}
	}()

	vars := mux.Vars(r)
	id := vars["id"]
	examID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return
	}

	res, err := db.DeleteRows(config.ScoreTable, config.ExamIdx, examID)
	if err != nil {
		log.Println(err)
		return
	}

	sendResponse(&models.GenericResponse{Message: fmt.Sprintf("Successfully deleted %v exams", res)}, http.StatusOK, w)
}

// GetAllExams gets a list of all examTestData that have been recorded (every record in data store)
func GetAllExams(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			SendGenericInternalServerError(w, r)
			return
		}
	}()

	res, err := db.GetRows(config.ScoreTable, config.IdFld)
	if err != nil {
		log.Println(err)
		return
	}

	response := &models.AllExamsListResponse{}
	for _, score := range res {
		response.Exams = append(response.Exams, score.(models.StudentExam))
	}

	sendResponse(response, http.StatusOK, w)
}

// GetAllUniqueExamIDs lists all the unique examTestData that have been recorded
func GetAllUniqueExamIDs(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			SendGenericInternalServerError(w, r)
			return
		}
	}()

	res, err := db.GetRows(config.ScoreTable, config.UniqueExamsIdx)
	if err != nil {
		log.Println(err)
		return
	}

	response := &models.AllUniqueExamsListResponse{}
	for _, score := range res {
		response.Exams = append(response.Exams, score.(models.StudentExam).Exam)
	}

	sendResponse(response, http.StatusOK, w)
}

// GetExamByID lists all the results for the specified exam, and provide the average score across all students
func GetExamByID(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			SendGenericInternalServerError(w, r)
			return
		}
	}()

	vars := mux.Vars(r)
	id := vars["id"]
	examID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return
	}

	response := &models.ExamByIDResponse{Exam: examID}
	res, err := db.GetRows(config.ScoreTable, config.ExamIdx, examID)
	if err != nil {
		log.Println(err)
		return
	}

	count := len(res)
	if count == 0 {
		SendGenericNotFoundResponse(w, r)
		return
	}

	var sum float64
	for _, score := range res {
		response.Scores = append(response.Scores, models.ExamScorePerStudent{Student: score.(models.StudentExam).StudentID, Score: score.(models.StudentExam).Score})
		sum += score.(models.StudentExam).Score
	}

	avg := sum / float64(count)
	response.Average = avg

	sendResponse(response, http.StatusOK, w)
}
