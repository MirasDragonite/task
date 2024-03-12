package services

import (
	"context"
	"miras/internal/controllers"
	"miras/internal/models"
	"net/http"
)

type Auth interface {
	Register(user models.Register) error
	Login(ctx context.Context, login models.Login) (*http.Cookie, error)
	Logout(ctx context.Context, cookie *http.Cookie) error
}

type Book interface {
	CreateBook(ctx context.Context, book models.Book) error
	GetBookByID(ctx context.Context, id int) (*models.Book, error)
	DeleteBook(ctx context.Context, id int) error
	GetAllBooks(ctx context.Context) ([]models.Book, error)
}

type Service struct {
	Auth
	Book
}

func NewService(repo *controllers.Repository) *Service {
	return &Service{Auth: newAuthService(repo), Book: newBookService(repo)}
}
