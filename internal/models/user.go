package models

import "time"

type User struct {
	ID       int
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	ID         int
	UserID     int
	Token      string
	ExpireDate time.Time
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
