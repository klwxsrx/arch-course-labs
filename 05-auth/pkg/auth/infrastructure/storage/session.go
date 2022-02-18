package storage

import (
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/auth/infrastructure/auth"
)

type sessionStorage struct {
	sessions map[string]*auth.Session
}

func (s *sessionStorage) Add(sessionID string, session *auth.Session) error {
	s.sessions[sessionID] = session
	return nil
}

func (s *sessionStorage) Get(sessionID string) (*auth.Session, error) {
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, auth.ErrSessionNotFound
	}
	return session, nil
}

func (s *sessionStorage) Remove(sessionID string) error {
	delete(s.sessions, sessionID)
	return nil
}

func NewSessionStorage() auth.SessionStorage {
	return &sessionStorage{make(map[string]*auth.Session)}
}
