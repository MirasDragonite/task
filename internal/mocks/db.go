package mocks

import (
	"database/sql"
	"fmt"
	"miras/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "mock.db")
	if err != nil {
		return nil, err
	}
	query := `
	DROP TABLE IF EXISTS books;
	DROP TABLE IF EXISTS sessions;
 	DROP TABLE IF EXISTS users;
	CREATE TABLE IF NOT  EXISTS users(id INTEGER PRIMARY KEY, username TEXT,email TEXT NOT NULL UNIQUE,hash_password TEXT NOT NULL,role TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS sessions(id INTEGER PRIMARY KEY,user_id INTEGER,token TEXT NOT NULL UNIQUE,expired_date TEXT NOT NULL,FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);	
	CREATE TABLE IF NOT EXISTS books(id INTEGER PRIMARY KEY,name TEXT NOT NULL,author TEXT NOT NULL,genre TEXT NOT NULL,year TEXT NOT NULL)
	`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to mock database")
	return db, nil
}

func CreateBook(db *sql.DB, book models.Book) error {
	query := `INSERT INTO books(name,author,genre,year) VALUES ($1,$2,$3,$4)`

	_, err := db.Exec(query, book.Name, book.Author, book.Genre, book.CreationYear)
	if err != nil {
		return err
	}

	return nil
}

func GetBookByID(db *sql.DB, id int) (models.Book, error) {
	var book models.Book
	query := `SELECT * FROM books WHERE id=$1`

	row := db.QueryRow(query, id)

	err := row.Scan(&book.Id, &book.Name, &book.Author, &book.Genre, &book.CreationYear)
	if err != nil {
		return models.Book{}, err
	}
	return book, nil
}

func GetAllBooks(db *sql.DB) ([]models.Book, error) {

	var books []models.Book

	query := `SELECT * FROM books`

	row, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var book models.Book
		err = row.Scan(&book.Id, &book.Name, &book.Author, &book.Genre, &book.CreationYear)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}
