package mocks

import (
	"database/sql"
	"miras/internal/models"
	"testing"
)

var mockDB *sql.DB
var err error
var expectedBook1 models.Book = models.Book{Id: 2, Name: "Indriver", Author: "Somedude", Genre: "dgs", CreationYear: "2020"}
var expectedBook2 models.Book = models.Book{Id: 1, Name: "CleanArchitecture", Author: "Somedude2", Genre: "it", CreationYear: "2023"}
var expctedBooks []models.Book = []models.Book{expectedBook2, expectedBook1}

func TestDBOpen(t *testing.T) {
	mockDB, err = NewDB()
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)
		return
	}
}
func TestBookFunctions(t *testing.T) {

	mockDB, err = NewDB()
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)

	}
	err = CreateBook(mockDB, expectedBook1)
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)
	}
}

func TestGetBookByID(t *testing.T) {
	mockDB, err = NewDB()
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)

	}
	err = CreateBook(mockDB, expectedBook2)
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)
	}
	book, err := GetBookByID(mockDB, expectedBook1.Id)
	if err == nil && book.Id != expectedBook2.Id {
		t.Errorf("Expected %d id, but got %d", expectedBook2.Id, book.Id)
	}

	book2, err := GetBookByID(mockDB, expectedBook2.Id)
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)
	}
	if book2.Id != expectedBook2.Id && book2.Name != expectedBook2.Name && book2.Author != expectedBook2.Author && book2.Genre != expectedBook2.Genre && book2.CreationYear != expectedBook2.CreationYear {

		t.Errorf("Expectd %v got %v", expectedBook2, book2)
	}

}

func TestGetAllBooks(t *testing.T) {

	mockDB, err = NewDB()
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)

	}
	err = CreateBook(mockDB, expectedBook2)
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)
	}
	err = CreateBook(mockDB, expectedBook1)
	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)
	}

	books, err := GetAllBooks(mockDB)

	if err != nil {
		t.Errorf("Expected %v got %v", nil, err)

	}
	for i := 0; i < len(books); i++ {
		if books[i].Id != expctedBooks[i].Id && books[i].Name != expctedBooks[i].Name && books[i].Author != expctedBooks[i].Author && books[i].Genre != expctedBooks[i].Genre && books[i].CreationYear != expctedBooks[i].CreationYear {
			t.Errorf("Expected %v got %v", books[i], expctedBooks[i])
		}
	}
}
