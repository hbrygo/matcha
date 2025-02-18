package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // Important: le _ est nécessaire
	"log"
	"os"
)

const dbPath = "./data/users.db"

func InitDB() (*sql.DB, error) {
	// verify is folder exists
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err := os.MkdirAll("./data", 0755)
		if err != nil {
			return nil, err
		}
	}

	// open connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// verify if tabes exist
	if !tablesExist(db) {
		err = createTables(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func tablesExist(db *sql.DB) bool {
	query := `
        SELECT name FROM sqlite_master 
        WHERE type='table' AND name='users'
    `
	var name string
	err := db.QueryRow(query).Scan(&name)
	return err == nil
}

func createTables(db *sql.DB) error {
	// read content of shema.sql file
	schema, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		return err
	}

	// execute schema.sql
	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	log.Println("Tables created successfully")
	return nil
}
