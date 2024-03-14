package controllers

import (
	"context"
	"database/sql"
	"miras/internal/models"
	"strings"
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

func (r *AuthRepo) CreateUser(user models.Register) (int64, error) {

	query := `INSERT INTO users(username,email,hash_password,role) VALUES($1,$2,$3,$4) RETURNING id;`
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *AuthRepo) SelectUser(login models.Login) (models.User, error) {

	var user models.User

	query := `SELECT * FROM users WHERE email=$1`
	row := r.db.QueryRow(query, login.Email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *AuthRepo) GetAllUserPermissions(ctx context.Context, userId int64) (map[string]bool, error) {

	var permissions = make(map[string]bool)

	query := `SELECT permissions.code
	FROM permissions
		INNER JOIN user_permissions ON user_permissions.permission_id = permissions.id
		INNER JOIN users ON user_permissions.user_id = users.id
	WHERE users.id = $1 `

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return permissions, err
	}

	for rows.Next() {
		var permission string
		err = rows.Scan(&permission)

		if err != nil {
			return nil, err
		}

		permissions[permission] = true
	}

	return permissions, nil
}

func (r *AuthRepo) AddForUser(userID int64, codes ...string) error {

	query := `INSERT INTO user_permissions(user_id,permission_id)
	SELECT $1,id FROM permissions WHERE code IN (?` + strings.Repeat(", ?", len(codes)-1) + `)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := make([]interface{}, len(codes)+1)
	args[0] = userID
	for i, code := range codes {
		args[i+1] = code
	}

	_, err = stmt.ExecContext(ctx, args...)
	return err
}
