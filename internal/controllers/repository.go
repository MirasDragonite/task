package controllers

import (
	"context"
	"database/sql"
	"miras/internal/models"

	"github.com/go-redis/cache/v9"
)

type Auth interface {
	CreateUser(user models.Register) error
	SelectUser(login models.Login) (models.User, error)
	GetSessionByUserID(id int) (models.Session, error)
	GetSessionByToken(token string) (models.Session, error)
	CreateSession(ctx context.Context, session models.Session) error
	DeleteToken(ctx context.Context, token string) error
}

type Book interface {
	CreateBook(ctx context.Context, book models.Book) error
	GetBookByID(ctx context.Context, id int) (*models.Book, error)
	GetAllBooks(ctx context.Context) ([]models.Book, error)
	DeleteBook(ctx context.Context, id int) error
}
type Repository struct {
	Auth
	Book
}

func NewRepository(db *sql.DB, cache *cache.Cache) *Repository {

	return &Repository{
		Auth: newAuthRepo(db, cache),
		Book: newBookRepo(db, cache),
	}
}
