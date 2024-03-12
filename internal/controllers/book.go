package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"miras/internal/models"
	"time"

	"github.com/go-redis/cache/v9"
)

type BookRepo struct {
	db    *sql.DB
	cache *cache.Cache
}

func newBookRepo(db *sql.DB, cache *cache.Cache) *BookRepo {
	return &BookRepo{db: db, cache: cache}
}

func (r *BookRepo) CreateBook(ctx context.Context, book models.Book) error {

	query := `INSERT INTO books(name,author,genre,year) VALUES ($1,$2,$3,$4)`

	_, err := r.db.Exec(query, book.Name, book.Author, book.Genre, book.CreationYear)
	if err != nil {
		return err
	}

	err = r.DeleteCacheFor(ctx, "all")
	if err != nil {
		return err
	}
	return nil
}

func (r *BookRepo) GetBookByID(ctx context.Context, id int) (*models.Book, error) {
	var book models.Book

	key := fmt.Sprintf("%d", id)
	err := r.cache.Get(ctx, key, &book)
	if err == nil {
		fmt.Println("from cache", book)
		return &book, nil
	}

	query := `SELECT * FROM books WHERE id=$1`

	row := r.db.QueryRow(query, id)

	err = row.Scan(&book.Id, &book.Name, &book.Author, &book.Genre, &book.CreationYear)
	if err != nil {
		return nil, err
	}

	err = r.cache.Set(&cache.Item{Ctx: ctx, Key: key, Value: book, TTL: time.Second * 40})
	if err != nil {
		return nil, err
	}
	fmt.Println("cached:", book)
	return &book, nil
}

func (r *BookRepo) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	var books []models.Book
	key := "all"

	err := r.cache.Get(ctx, key, &books)
	if err == nil {
		fmt.Println("Get data from cache")
		return books, nil
	}

	query := `SELECT * FROM books`

	row, err := r.db.Query(query)
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
	err = r.cache.Set(&cache.Item{Ctx: ctx, Key: key, Value: books, TTL: time.Minute * 3})
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepo) DeleteBook(ctx context.Context, id int) error {

	query := "DELETE FROM books where id=$1"

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	err = r.DeleteCacheFor(ctx, "all", fmt.Sprintf("%d", id))
	if err != nil {
		return err
	}

	return nil
}

func (r *BookRepo) DeleteCacheFor(ctx context.Context, s ...string) error {
	for _, v := range s {
		err := r.cache.Delete(ctx, v)
		if err != nil {
			return err
		}
	}
	return nil
}
