package controllers

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		return nil, err
	}
	query := `
	DROP TABLE IF EXISTS books;
	DROP TABLE IF EXISTS sessions;
 	DROP TABLE IF EXISTS users;
	CREATE TABLE IF NOT  EXISTS users(id INTEGER PRIMARY KEY, username TEXT,email TEXT NOT NULL UNIQUE,hash_password TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS sessions(id INTEGER PRIMARY KEY,user_id INTEGER,token TEXT NOT NULL UNIQUE,expired_date TEXT NOT NULL,FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);	
	CREATE TABLE IF NOT EXISTS books(id INTEGER PRIMARY KEY,name TEXT NOT NULL,author TEXT NOT NULL,genre TEXT NOT NULL,year TEXT NOT NULL)
	`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database")
	return db, nil
}
