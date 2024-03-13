package rest

import "github.com/go-chi/chi/v5"

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/sign-up", h.signUp)
	r.Post("/sign-in", h.signIN)
	r.Post("/log-out", h.RequireAuth(h.logout))
	r.Post("/create-book", h.RequireAuth(h.RequirePermissions(h.CreateBook, "book:create")))
	r.Post("/delete-book", h.RequireAuth(h.RequirePermissions(h.deleteBook, "book:delete")))
	r.Get("/book", h.RequireAuth(h.RequirePermissions(h.getBook, "book:read")))
	r.Get("/get-books", h.RequireAuth(h.RequirePermissions(h.getAllBooks, "book:read_all")))
	return r
}
