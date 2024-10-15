package middleware

import (
	"context"
	"fmt"
	"net/http"
	"steganography/internal/auth"
)

type FavContextKey string

var KEY = FavContextKey("userID")

// Make the session manager accessible within the package
var sessionManager *auth.SessionManager

func InitSessionManager(sm *auth.SessionManager) {
	sessionManager = sm
}

// add middleware to check if the session is valid
func ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("Error retrieving cookie:", err)
			return
		}

		session, exists := sessionManager.GetSession(sessionCookie.Value)
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Store the userID in the context
		ctx := context.WithValue(r.Context(), KEY, session.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r) // start the next handler
	})
}
