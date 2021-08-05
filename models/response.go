package models

// GenericResponse is the generic model for providing a response to the client
type GenericResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// AllStudentListResponse defines the response returned when retrieving a list of students
type AllStudentListResponse struct {
	Students []string `json:"students"`
}

// StudentByIDResponse defines the response returned when retrieving a specific student record
type StudentByIDResponse struct {
	Student string              `json:"student"`
	Exams   []StudentExamScores `json:"exams"`
	Average float64             `json:"average"`
}

// StudentExamScores is a simple struct that contains an exam id and score
type StudentExamScores struct {
	Exam  int     `json:"exam"`
	Score float64 `json:"score"`
}

// AllUniqueExamsListResponse is the response returned when retrieving a list of all unique exams
type AllUniqueExamsListResponse struct {
	Exams []int `json:"exams"`
}

// AllExamsListResponse is the response returned when retrieving a list of all exams
type AllExamsListResponse struct {
	Exams []StudentExam `json:"exams"`
}

// ExamByIDResponse is the response returned when retrieving a specific exam record
type ExamByIDResponse struct {
	Exam    int                   `json:"exam"`
	Scores  []ExamScorePerStudent `json:"scores"`
	Average float64               `json:"average"`
}

// ExamScorePerStudent is a simple struct that contains a student id and score
type ExamScorePerStudent struct {
	Student string  `json:"student"`
	Score   float64 `json:"score"`
}
