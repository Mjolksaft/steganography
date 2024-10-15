package auth

import (
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
	return sessionID
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	return session, exists
}

func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, sessionID)
}

func (sm *SessionManager) CleanupExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for id, session := range sm.sessions {
		if time.Now().After(session.ExpiresAt) {
			delete(sm.sessions, id)
		}
	}
}
