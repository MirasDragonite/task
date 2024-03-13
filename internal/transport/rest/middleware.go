package rest

import (
	"net/http"
	"time"
)

func (h *Handler) RequireAuth(request http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//take current session from our user
		cookie, err := r.Cookie("Token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var session http.Cookie

		//get session from cache and compare it with current session
		err = h.cache.Get(r.Context(), "session", &session)
		if err != nil || session.Value != cookie.Value || session.Expires.Before(time.Now()) {

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
