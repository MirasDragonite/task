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
	Logout(cookie *http.Cookie)
	GetAllUserPermissions(ctx context.Context, userID int64) (map[string]bool, error)
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
