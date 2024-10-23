package middleware

import (
	"context"
	"log/slog"
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

// middleware to check if the session is valid
func ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			slog.Warn("Failed to retrieve session cookie", slog.String("err", err.Error()))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session, exists := sessionManager.GetSession(sessionCookie.Value)
		if !exists {
			slog.Warn("Session does not exist", slog.String("sessionID", sessionCookie.Value))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Store the userID in the context
		ctx := context.WithValue(r.Context(), KEY, session.UserID)
		r = r.WithContext(ctx)

		slog.Info("Session validated", slog.String("sessionID", sessionCookie.Value), slog.String("userID", session.UserID))
		next.ServeHTTP(w, r) // start the next handler
	})
}

// ValidateAdmin checks if the user has admin privileges
func ValidateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* TODO: Check is is_admin is true */
		next.ServeHTTP(w, r) // start the next handler
	})
}
