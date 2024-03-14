package rest

import (
	"encoding/json"
	"log"
	"miras/internal/models"
	"net/http"
	"time"

	"github.com/go-redis/cache/v9"
)

// Handler to register new user
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	var user models.Register

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatal("Error druing decode")
		return
	}

	err = h.Service.Auth.Register(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("User successfully create"))
}

// handler to login into the sytem
func (h *Handler) signIN(w http.ResponseWriter, r *http.Request) {

	var user models.Login

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// we use context
	cookie, err := h.Service.Login(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = h.cache.Set(&cache.Item{Ctx: r.Context(), Key: "session", Value: cookie, TTL: time.Minute * 10})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(200)
	w.Write([]byte("User successfully create"))
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	// don't take err cause this function wrapped with middleware, and it's already check cookies
	cookie, _ := r.Cookie("Token")

	h.Service.Auth.Logout(cookie)

	err := h.cache.Delete(r.Context(), "session")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write([]byte("successfully logout"))
}
