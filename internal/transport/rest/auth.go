package rest

import (
	"encoding/json"
	"log"
	"miras/internal/models"
	"net/http"
)

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
	}

	w.WriteHeader(200)
	w.Write([]byte("User successfully create"))
}

func (h *Handler) signIN(w http.ResponseWriter, r *http.Request) {

	var user models.Login

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	cookie, err := h.Service.Login(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(200)
	w.Write([]byte("User successfully create"))
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("Token")
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	err = h.Service.Auth.Logout(r.Context(), cookie)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("successfully logout"))
}
