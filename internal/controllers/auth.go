package controllers

import (
	"database/sql"
	"fmt"
	"miras/internal/models"
	"time"

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
	fmt.Println("GIG1")
	return user, nil
}

func (r *AuthRepo) CreateSession(session models.Session) error {

	query := `INSERT INTO sessions(user_id,token,expired_date) VALUES($1,$2,$3)`
	expDate := session.ExpireDate.Format("2006-01-02 15:04:05")
	_, err := r.db.Exec(query, session.UserID, session.Token, expDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) GetSessionByToken(token string) (models.Session, error) {
	var session models.Session
	query := `SELECT * FROM sessions WHERE token=$1`

	row := r.db.QueryRow(query, token)
	var date string
	err := row.Scan(&session.ID, &session.UserID, &session.Token, &date)
	if err != nil {
		return models.Session{}, err
	}
	session.ExpireDate, err = time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (r *AuthRepo) GetSessionByUserID(id int) (models.Session, error) {
	var session models.Session
	query := `SELECT * FROM sessions WHERE user_id=$1`

	row := r.db.QueryRow(query, id)
	var date string
	err := row.Scan(&session.ID, &session.UserID, &session.Token, &date)
	if err != nil {
		return models.Session{}, err
	}
	session.ExpireDate, err = time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}
