package config

import "github.com/hashicorp/go-memdb"

type Config struct {
	MemDBSchema  *memdb.DBSchema
	SSEServerUrl string
	PORT string
}

const EnvURL = "SSE_SERVER_URL"
const EnvPort = "APPLICATION_PORT"

// Define the table name, fields, and indexes for the in-memory data store
const (
	ScoreTable        = "score"
	StudentIdx        = "student_idx"
	ExamIdx           = "exam_idx"
	UniqueStudentsIdx = "u_student_idx"
	UniqueExamsIdx    = "u_exam_idx"
	IdFld             = "id"
	ExamFld           = "Exam"
	StudentFld        = "StudentID"
)

// DBSchema Define the schema used for the scores in-memory database
var DBSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		ScoreTable: {
			Name: ScoreTable,
			Indexes: map[string]*memdb.IndexSchema{
				IdFld: {
					Name:   IdFld,
					Unique: true,
					Indexer: &memdb.CompoundIndex{
						Indexes: []memdb.Indexer{
							&memdb.IntFieldIndex{Field: ExamFld},
							&memdb.StringFieldIndex{Field: StudentFld},
						},
						AllowMissing: true,
					},
				},
				StudentIdx: {
					Name:    StudentIdx,
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: StudentFld},
				},
				UniqueStudentsIdx: {
					Name:    UniqueStudentsIdx,
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: StudentFld},
				},
				ExamIdx: {
					Name:    ExamIdx,
					Unique:  false,
					Indexer: &memdb.IntFieldIndex{Field: ExamFld},
				},
				UniqueExamsIdx: {
					Name:    UniqueExamsIdx,
					Unique:  true,
					Indexer: &memdb.IntFieldIndex{Field: ExamFld},
				},
			},
		},
	},
}
