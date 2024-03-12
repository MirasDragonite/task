package rest

import "github.com/go-chi/chi/v5"

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/sign-up", h.signUp)
	r.Post("/sign-in", h.signIN)
	r.Post("/create-book", h.CreateBook)
	r.Post("/delete-book", h.deleteBook)
	r.Get("/book", h.getBook)
	r.Get("/get-books", h.getAllBooks)
	return r
}
