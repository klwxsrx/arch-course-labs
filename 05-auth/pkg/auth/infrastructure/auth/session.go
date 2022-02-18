package auth

import (
	"errors"
	"github.com/google/uuid"
)

type Session struct {
	UserID uuid.UUID
	Login  string
}

var ErrSessionNotFound = errors.New("session not found")

type SessionStorage interface {
	Add(sessionID string, session *Session) error
	Get(sessionID string) (*Session, error)
	Remove(sessionID string) error
}
