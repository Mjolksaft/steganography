package auth

import (
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) CreateSession(userID string, sessionExpiration time.Time) string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionID := uuid.NewString()
	session := &Session{
		UserID: userID,
	}
	sm.sessions[sessionID] = session

	slog.Info("Session created", slog.String("sessionID", sessionID), slog.String("userID", userID), slog.Time("expiresAt", sessionExpiration))
	return sessionID
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if exists {
		slog.Info("Session retrieved", slog.String("sessionID", sessionID), slog.String("userID", session.UserID))
		return session, exists
	}

	return nil, exists
}

func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.sessions[sessionID]; exists {
		delete(sm.sessions, sessionID)
		slog.Info("Session deleted", slog.String("sessionID", sessionID))
	}
}

func (sm *SessionManager) CleanupExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for id, session := range sm.sessions {
		if time.Now().After(session.ExpiresAt) {
			delete(sm.sessions, id)
			slog.Info("Expired session cleaned up", slog.String("sessionID", id), slog.String("userID", session.UserID))
		}
	}
}
