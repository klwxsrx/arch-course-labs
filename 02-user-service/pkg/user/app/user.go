package app

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists                   = errors.New("user already exists")
	ErrUserWithSpecifiedEmailAlreadyExists = errors.New("user with specified email already exists")
	ErrUserNotFound                        = errors.New("user not found")
)

type UserRepository interface {
	Create(u *User) error
	Get(id uuid.UUID) (*User, error)
	Update(u *User) error
	Delete(id uuid.UUID) error
}

type User struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
}
