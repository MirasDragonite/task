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

	return s.repo.Auth.CreateUser(user)

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

	newToken, err := GenerateJWTToken()

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
