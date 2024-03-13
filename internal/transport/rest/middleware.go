package rest

import (
	"fmt"
	"miras/internal/models"
	"net/http"
	"time"
)

func (h *Handler) RequirePermissions(request http.HandlerFunc, permissions ...string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var session models.Session
		err := h.cache.Get(r.Context(), "session", &session)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(session)
		availablePermissions, err := h.Service.Auth.GetAllUserPermissions(r.Context(), int64(session.UserID))
		if err != nil {
			http.Error(w, "bad request"+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(permissions, availablePermissions)
		if !Includes(permissions, availablePermissions.Permissions) {
			http.Error(w, "Access denied", http.StatusBadRequest)
			return
		}

		request.ServeHTTP(w, r)
	})
}

func (h *Handler) RequireAuth(request http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//take current session from our user
		cookie, err := r.Cookie("Token")
		if err != nil {

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var session models.Session

		//get session from cache and compare it with current session
		err = h.cache.Get(r.Context(), "session", &session)
		if err != nil || session.Token != cookie.Value || session.ExpireDate.Before(time.Now()) {

			// if they don't match  then delete the token from our cache and show error with code 401
			err = h.cache.Delete(r.Context(), "session")
			if err != nil {
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// let use the handle func if everything is ok
		request.ServeHTTP(w, r)
	})
}

func Includes(permissions []string, available map[string]bool) bool {

	for _, ch := range permissions {
		if !available[ch] {
			return false
		}
	}
	return true
}
