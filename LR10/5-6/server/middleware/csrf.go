package middleware

import (
	"net/http"
	"server/controllers"
)

func CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			token := r.Header.Get("X-CSRF-Token")
			if token == "" {
				http.Error(w, "CSRF token missing", http.StatusForbidden)
				return
			}
			err := controllers.ValidateCSRFToken(token, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
