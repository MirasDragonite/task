package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func (h *Handler) RequirePermissions(request http.HandlerFunc, permissions ...string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var session http.Cookie
		// get session from cache to get from there user permissions
		err := h.cache.Get(r.Context(), "session", &session)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims, err := ParseToken(session.Value)
		if err != nil {
			http.Error(w, "Bad request:"+err.Error(), http.StatusBadRequest)
			return
		}
		permissionsRaw, ok := claims["permissions"]
		//if permissions not found in claims"
		if !ok {
			http.Error(w, "Bad request:", http.StatusBadRequest)
			return
		}

		availablePermissions, ok := permissionsRaw.(map[string]interface{})
		//check if permission is type of map[string]interface
		if !ok {
			http.Error(w, "Bad request:", http.StatusBadRequest)
			return
		}
		//chech if user have required permissions
		if !Includes(permissions, availablePermissions) {
			http.Error(w, "Access denied", http.StatusBadRequest)
			return
		}

		// let use the handle func if everything is ok
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

// func to check is this user have required permissions
func Includes(permissions []string, available map[string]interface{}) bool {

	for _, ch := range permissions {
		_, ok := available[ch].(bool)
		//check if the value exist in map
		if !ok {
			return false
		}
	}

	return true
}

func ParseToken(signedToken string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("qwertyacid12345acidqwerty"), nil
	})
	if err != nil {
		return nil, err
	}

	//if token is valid
	if parsedToken.Valid {
		claims := parsedToken.Claims.(jwt.MapClaims)
		return claims, nil
	}

	return nil, errors.New("Token is invalid")

}
