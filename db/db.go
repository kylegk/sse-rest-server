package db

import (
	"log"

	"github.com/hashicorp/go-memdb"
)

var db *memdb.MemDB

// InitDB initializes the datastore
func InitDB(schema *memdb.DBSchema) error {
	if schema == nil {
		panic("cannot initialize database: missing schema")
	}

	conn, err := memdb.NewMemDB(schema)
	if err != nil {
		panic("unable to initialize database: " + err.Error())
	}

	db = conn

	return nil
}

// UpsertRow inserts a row into the database if it doesn't exist, or updates the existing value(s) if it does
func UpsertRow(table string, record interface{}) error {
	if db == nil {
		panic("database connection has not been initialized")
	}

	txn := db.Txn(true)
	defer txn.Abort()

	err := txn.Insert(table, record)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

// DeleteRow removes a single row from the database
// Not currently used
func DeleteRow(table string, record interface{}) error {
	if db == nil {
		panic("database connection has not been initialized")
	}

	txn := db.Txn(true)
	defer txn.Abort()

	err := txn.Delete(table, record)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

// DeleteRows removes multiple rows from the database
func DeleteRows(table string, idx string, args ...interface{}) (int, error) {
	if db == nil {
		panic("database connection has not been initialized")
	}

	txn := db.Txn(true)
	defer txn.Abort()

	count, err := txn.DeleteAll(table, idx, args...)
	if err != nil {
		return 0, err
	}

	log.Println("delete rows: ", count)

	txn.Commit()

	return count, nil
}

// GetRows retrieves all the rows in the database
func GetRows(table string, idx string, args ...interface{}) ([]interface{}, error) {
	if db == nil {
		panic("database connection has not been initialized")
	}

	txn := db.Txn(false)

	it, err := txn.Get(table, idx, args...)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		results = append(results, obj)
	}

	return results, nil
}
