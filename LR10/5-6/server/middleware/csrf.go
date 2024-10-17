package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("your-secret-key"))

func CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			session, err := Store.Get(r, "session")
			if err != nil {
				http.Error(w, "сессия не найдена", http.StatusUnauthorized)
				return
			}
			csrfToken := session.Values["csrf_token"]
			if csrfToken == nil || csrfToken != r.Header.Get("X-CSRF-Token") {
				http.Error(w, "недействительный CSRF-токен", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
