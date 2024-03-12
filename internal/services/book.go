package services

import (
	"context"
	"miras/internal/controllers"
	"miras/internal/models"
)

type BookService struct {
	repo *controllers.Repository
}

func newBookService(repo *controllers.Repository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(ctx context.Context, book models.Book) error {

	return s.repo.CreateBook(ctx, book)
}

func (s *BookService) GetBookByID(ctx context.Context, id int) (*models.Book, error) {
	return s.repo.Book.GetBookByID(ctx, id)
}
func (s *BookService) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return s.repo.Book.GetAllBooks(ctx)
}
func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.repo.Book.DeleteBook(ctx, id)
}
