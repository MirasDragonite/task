package rest

import "github.com/go-chi/chi/v5"

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/sign-up", h.signUp)
	r.Post("/sign-in", h.signIN)
	r.Post("/log-out", h.RequireAuth(h.logout))
	r.Post("/create-book", h.RequireAuth(h.CreateBook))
	r.Post("/delete-book", h.RequireAuth(h.deleteBook))
	r.Get("/book", h.RequireAuth(h.getBook))
	r.Get("/get-books", h.RequireAuth(h.getAllBooks))
	return r
}
