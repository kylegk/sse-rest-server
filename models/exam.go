package models

// StudentExam defines event messages returned from the sse client
type StudentExam struct {
	Exam      int     `json:"exam"`
	StudentID string  `json:"studentid"`
	Score     float64 `json:"score"`
}
