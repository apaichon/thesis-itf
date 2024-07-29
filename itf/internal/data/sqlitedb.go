package data

import (
	"database/sql"

	"fmt"
	"log"
	"sync"

	config "github.com/apaichon/thesis-itf/itf/config"
	_ "github.com/mattn/go-sqlite3"
)

// SqliteDB represents the SQLite database
type SqliteDB struct {
	Connection *sql.DB
}

var instance *SqliteDB
var once sync.Once

// NewSqliteDB initializes a new instance of the SqliteDB struct
func NewSqliteDB() *SqliteDB {
	once.Do(func() {
		config := config.NewConfig()
		conn, err := sql.Open("sqlite3", config.TempDBPath)
		if err != nil {
			log.Fatal(err)
		}
		instance = &SqliteDB{conn}
	})
	return instance
}

func (db *SqliteDB) Open() error {
	if db.Connection == nil {
		config := config.NewConfig()
		conn, err := sql.Open("sqlite3", config.TempDBPath)
		if err != nil {
			return err
		}
		instance = &SqliteDB{conn}
	}
	return nil
}

// Close closes the database connection
func (db *SqliteDB) Close() error {
	if db.Connection == nil {
		return nil
	}
	return db.Connection.Close()
}

// Insert inserts data into the specified table
func (db *SqliteDB) Insert(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, nil
}

// Query executes a query and returns rows
func (db *SqliteDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return rows, nil
}

// QueryRow executes a query that is expected to return at most one row
func (db *SqliteDB) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	row := db.Connection.QueryRow(query, args...)
	return row, nil
}

// Delete executes a delete statement
func (db *SqliteDB) Delete(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, nil
}

// Update executes an update statement
func (db *SqliteDB) Update(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, nil
}

func (db *SqliteDB) Begin() (*sql.Tx, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	return tx, nil
}

func (db *SqliteDB) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	return stmt, nil
}

func (db *SqliteDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Connection.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}
	return result, nil
}
