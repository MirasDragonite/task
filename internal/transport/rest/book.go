package rest

import (
	"encoding/json"
	"miras/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {

	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	err = h.Service.Book.CreateBook(r.Context(), book)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	w.WriteHeader(200)
	w.Write([]byte("Book successfully created"))
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	book, err := h.Service.Book.GetBookByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
}
func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request) {

	books, err := h.Service.Book.GetAllBooks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	err = h.Service.Book.DeleteBook(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	w.WriteHeader(200)
	w.Write([]byte("Book successfully delete"))
}
