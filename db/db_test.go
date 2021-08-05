package db

import (
	"github.com/hashicorp/go-memdb"
	"github.com/kylegk/sse-rest-server/config"
	"github.com/kylegk/sse-rest-server/models"
	"log"
	"testing"
)

var validSchema = config.DBSchema
var validTable = config.ScoreTable
var validIdx = config.ExamIdx

// TestDBInitNil validates the database init will fail when passed a nil value
func TestDBInitNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	err := InitDB(nil)
	if err != nil {
		t.Error(err.Error())
	}
}

// TestDBInitSchemaError validates the database init will fail when passed a malformed schema
func TestDBInitSchemaError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	var invalidSchema = &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"foo": {
				Name: "bar",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id_",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "id"},
					},
				},
			},
		},
	}

	err := InitDB(invalidSchema)
	if err != nil {
		t.Error(err.Error())
	}
}

// TestDeleteRowsUninitializedDatabase validates behavior when attempting to call DeleteRows when no database connection exists
func TestDeleteRowsUninitializedDatabase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	_, err := DeleteRows(validTable, config.ExamIdx, 1)
	if err != nil {
		log.Println(err)
	}
}

// TestValidDBInit validates that the database will initialize with a valid schema
func TestValidDBInit(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}
}

// TestUpsertRowUninitializedDatabase validates the proper panic behavior when trying to upsert on an uninitialized database
func TestUpsertRowUninitializedDatabase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	err := InitDB(nil)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	record := models.StudentExam{}
	err = UpsertRow("table", record)
	if err != nil {
		log.Println(err)
	}
}

// TestUpsertRowInvalidTable validates the proper error handling when attempting to upsert into a table that does not exist in the database
func TestUpsertRowInvalidTable(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	record := models.StudentExam{
		Exam:      111,
		StudentID: "test",
		Score:     100,
	}

	err = UpsertRow("table", record)
	if err == nil {
		t.Errorf("The insert should have failed")
	}
}

// TestUpsertRowInvalidRecordType validates that insert fails when the record has a mismatched data type
func TestUpsertRowInvalidRecordType(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	type Foo struct {
		Bar string
	}

	record := Foo{Bar: "fail"}

	err = UpsertRow(validTable, record)
	if err == nil {
		t.Errorf("The insert should have failed")
	}
}

// TestUpsertRow validates that a row is inserted into the database table
func TestUpsertRow(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	record := models.StudentExam{
		Exam:      111,
		StudentID: "test",
		Score:     100,
	}

	err = UpsertRow(validTable, record)
	if err != nil {
		t.Errorf("The insert should have succeeded")
	}
}

// TestDeleteRowsInvalidTable validates delete will fail when provided an invalid table
func TestDeleteRowsInvalidTable(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	_, err = DeleteRows("foo", "bar", 12345)
	if err == nil {
		t.Errorf("The delete should have failed")
	}
}

// TestDeleteRowsInvalidTable validates delete will fail when provided an invalid index
func TestDeleteRowsInvalidIndex(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	_, err = DeleteRows(validTable, "bar", 12345)
	if err == nil {
		t.Errorf("The delete should have failed")
	}
}

// TestDeleteRows verifies inserted rows will be deleted from the database
func TestDeleteRows(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	record := models.StudentExam{
		Exam:      111,
		StudentID: "test",
		Score:     100,
	}
	err = UpsertRow(validTable, record)
	if err != nil {
		t.Errorf("Failed to insert prior to delete")
	}

	want := 1
	have, err := DeleteRows(validTable, validIdx, 111)
	if have != want {
		t.Errorf("Failed to delete the correct number of records; have: %v, want %v", have, want)
	}
}

// TestGetRowsInvalidTable validates the behavior of GetRows when provided an invalid table name
func TestGetRowsInvalidTable(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	_, err = GetRows("foo", validIdx, 12345)
	if err == nil {
		t.Errorf("The lookup should have failed")
	}
}

// TestGetRowsInvalidTable validates the behavior of GetRows when provided an invalid index
func TestGetRowsInvalidIndex(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	_, err = GetRows(validTable, "foo", 12345)
	if err == nil {
		t.Errorf("The lookup should have failed")
	}
}

// TestGetRows validates that GetRows works as expected
func TestGetRows(t *testing.T) {
	err := InitDB(validSchema)
	if err != nil {
		t.Errorf("The database failed to initialize")
	}

	record := models.StudentExam{
		Exam:      111,
		StudentID: "test",
		Score:     100,
	}
	err = UpsertRow(validTable, record)
	if err != nil {
		t.Errorf("Failed to insert prior to lookup")
	}

	rows, err := GetRows(validTable, validIdx, 111)
	want := 1
	have := len(rows)
	if have != want {
		t.Errorf("Failed to retrieve the correct number of records; have: %v, want %v", have, want)
	}
}
