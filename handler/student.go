package handler

import (
	"github.com/kylegk/sse-rest-server/config"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kylegk/sse-rest-server/db"
	"github.com/kylegk/sse-rest-server/models"
)

// GetAllStudents lists all students that have received at least one test score
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			SendGenericInternalServerError(w, r)
			return
		}
	}()

	res, err := db.GetRows(config.ScoreTable, config.UniqueStudentsIdx)
	if err != nil {
		log.Println(err)
		return
	}

	response := &models.AllStudentListResponse{}
	for _, score := range res {
		response.Students = append(response.Students, score.(models.StudentExam).StudentID)
	}

	sendResponse(response, http.StatusOK, w)
}

// GetStudentByID lists the exam results for the specified student, and provides the student's average score across all examTestData
func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			SendGenericInternalServerError(w, r)
			return
		}
	}()

	vars := mux.Vars(r)
	studentID := vars["id"]
	response := &models.StudentByIDResponse{Student: studentID}

	res, err := db.GetRows(config.ScoreTable, config.StudentIdx, studentID)
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
		response.Exams = append(response.Exams, models.StudentExamScores{Exam: score.(models.StudentExam).Exam, Score: score.(models.StudentExam).Score})
		sum += score.(models.StudentExam).Score
	}

	avg := sum / float64(count)
	response.Average = avg

	sendResponse(response, http.StatusOK, w)
}
