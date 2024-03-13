package controllers

import (
	"database/sql"
	"miras/internal/models"

	"github.com/go-redis/cache/v9"
)

type AuthRepo struct {
	db    *sql.DB
	cache *cache.Cache
}

func newAuthRepo(db *sql.DB, cache *cache.Cache) *AuthRepo {
	return &AuthRepo{db: db, cache: cache}
}

func (r *AuthRepo) CreateUser(user models.Register) error {

	query := `INSERT INTO users(username,email,hash_password) VALUES($1,$2,$3)`

	_, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepo) SelectUser(login models.Login) (models.User, error) {

	var user models.User
	query := `SELECT * FROM users WHERE email=$1`

	row := r.db.QueryRow(query, login.Email)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
