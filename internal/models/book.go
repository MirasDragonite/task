package models

type Book struct {
	Id           int
	Name         string `json:"name"`
	Author       string `json:"author"`
	Genre        string `json:"genre"`
	CreationYear string `json:"year"`
}
