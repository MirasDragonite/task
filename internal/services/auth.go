package services

import (
	"context"
	"errors"
	"miras/internal/controllers"
	"miras/internal/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *controllers.Repository
}

func newAuthService(repo *controllers.Repository) *AuthService {
	return &AuthService{repo: repo}
}

// registration in service layer
func (s *AuthService) Register(user models.Register) error {
	hashPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPassword

	id, err := s.repo.Auth.CreateUser(user)
	if err != nil {
		return err
	}
	// set user permissions by his role
	permissions := []string{"book:read", "book:create"}
	if user.Role == "admin" {
		permissions = append(permissions, "book:delete", "book:read_all")
	} else if user.Role == "employee" {
		permissions = append(permissions, "book:read_all")
	}
	err = s.repo.Auth.AddForUser(id, permissions...)
	if err != nil {
		return err
	}
	return nil

}

// logining in service layer
func (s *AuthService) Login(ctx context.Context, login models.Login) (*http.Cookie, error) {
	// new cookie to define our session
	cookie := &http.Cookie{
		Name:     "Token",
		Path:     "/",
		HttpOnly: true,
	}

	user, err := s.repo.Auth.SelectUser(login)
	if err != nil {
		return nil, err
	}
	// comapare our passwords to check are they match
	if !doPasswordsMatch(user.Password, login.Password) {
		return nil, errors.New("password don't match")
	}

	// get user all permissions
	permissions, err := s.repo.GetAllUserPermissions(ctx, int64(user.ID))
	if err != nil {
		return nil, err
	}
	// create new token
	newToken, err := GenerateJWTToken(int64(user.ID), permissions)
	if err != nil {
		return nil, err
	}
	cookie.Value = newToken
	cookie.Expires = time.Now().Add(10 * time.Minute)

	return cookie, nil
}

func (s *AuthService) Logout(cookie *http.Cookie) {
	// values to delete our cookie
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	cookie.HttpOnly = false

}

func (s *AuthService) GetAllUserPermissions(ctx context.Context, userID int64) (map[string]bool, error) {
	return s.repo.Auth.GetAllUserPermissions(ctx, userID)
}

// function to hash password
func hashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// function to comapre passwords
func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func GenerateJWTToken(userID int64, permissions map[string]bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["permissions"] = permissions

	tokenString, err := token.SignedString([]byte("qwertyacid12345acidqwerty"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
