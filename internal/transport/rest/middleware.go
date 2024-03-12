package rest

import (
	"miras/internal/models"
	"net/http"
	"time"
)

func (h *Handler) RequireAuth(request http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("Token")
		if err != nil {
			http.Error(w, "UnAuthorized", http.StatusUnauthorized)
			return
		}
		var session models.Session
		err = h.cache.Get(r.Context(), "session", &session)
		if err != nil || session.Token != cookie.Value || session.ExpireDate.Before(time.Now()) {
			err = h.cache.Delete(r.Context(), "session")
			if err!=nil{
				return 
			}
			http.Error(w, "UnAuthorized", http.StatusUnauthorized)
			return
		}

		request.ServeHTTP(w, r)
	})
}
