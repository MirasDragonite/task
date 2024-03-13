package services

import (
	"context"
	"errors"
	"fmt"
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
	err = s.repo.Auth.AddForUser(id, "book:read", "book:create")
	if err != nil {
		return err
	}
	return nil

}

// logining in service layer
func (s *AuthService) Login(ctx context.Context, login models.Login) (*http.Cookie, models.Session, error) {
	// new cookie to define our session
	cookie := &http.Cookie{
		Name:     "Token",
		Path:     "/",
		HttpOnly: true,
	}

	user, err := s.repo.Auth.SelectUser(login)
	if err != nil {
		return nil, models.Session{}, err
	}
	// comapare our passwords to check are they match
	if !doPasswordsMatch(user.Password, login.Password) {
		return nil, models.Session{}, errors.New("password don't match")
	}

	newToken, err := GenerateJWTToken()

	if err != nil {
		return nil, models.Session{}, err
	}
	var session models.Session

	cookie.Value = newToken
	cookie.Expires = time.Now().Add(10 * time.Minute)
	session.Token = newToken
	session.ExpireDate = time.Now().Add(10 * time.Minute)
	session.UserID = user.ID
	return cookie, session, nil
}

func (s *AuthService) Logout(cookie *http.Cookie) {
	// values to delete our cookie
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	cookie.HttpOnly = false

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

// new token for our session
func GenerateJWTToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "@MIras"
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()

	tokenString, err := token.SignedString([]byte("qwertyacid12345acidqwerty"))

	if err != nil {
		fmt.Println("generating JWT Token failed")
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) GetAllUserPermissions(ctx context.Context, userID int64) (models.Permissions, error) {
	return s.repo.Auth.GetAllUserPermissions(ctx, userID)
}
