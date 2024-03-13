package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {
	//opening new connection to db
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		return nil, err
	}

	// taking our queries from .txt file, cause we don't have actualy migrations in sqlite, so this alternative.
	sqlBytes, err := os.ReadFile("migrations.txt")
	if err != nil {
		return nil, errors.New("failed to read file: " + err.Error())
	}
	query := string(sqlBytes)

	//inserting our data into db
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database")
	return db, nil
}
