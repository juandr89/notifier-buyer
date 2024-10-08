package middleware

import (
	"net/http"
)

func ApiKeyMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("x-api-key")
			if authHeader != apiKey {
				http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
