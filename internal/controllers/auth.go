package controllers

import (
	"context"
	"database/sql"
	"fmt"
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

	query := `INSERT INTO users(username,email,hash_password) VALUES($1,$2,$3) RETURNING id;`
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password)
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
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *AuthRepo) GetAllUserPermissions(ctx context.Context, userId int64) (models.Permissions, error) {
	fmt.Println("HERE I AM")
	var permissions models.Permissions

	permissions.Permissions = make(map[string]bool)

	err := r.cache.Get(ctx, "permissions", &permissions)
	if err == nil {
		fmt.Println("FROM THIS CACHE")
		return permissions, nil
	}

	query := `SELECT permissions.code
	FROM permissions
		INNER JOIN user_permissions ON user_permissions.permission_id = permissions.id
		INNER JOIN users ON user_permissions.user_id = users.id
	WHERE users.id = $1 `

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return permissions, err
	}
	permissions.UserID = userId
	for rows.Next() {
		var permission string
		err = rows.Scan(&permission)
		fmt.Println(permission)
		if err != nil {
			return models.Permissions{}, err
		}

		permissions.Permissions[permission] = true
	}

	err = r.cache.Set(&cache.Item{Ctx: ctx, Key: "permissions", Value: permissions, TTL: time.Minute * 15})
	if err != nil {
		return models.Permissions{}, err
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
