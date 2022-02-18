package app

import (
	"errors"
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Login           string
	EncodedPassword string
}

var ErrUserByLoginAlreadyExists = errors.New("user with specified login exists")
var ErrUserByLoginNotExists = errors.New("user with specified login doesn't exist")

type UserRepository interface {
	NextID() uuid.UUID
	GetByLogin(login string) (*User, error)
	Store(user *User) error
}

type UserService struct {
	repo       UserRepository
	pwdEncoder PasswordEncoder
}

func (s *UserService) Register(login, password string) error {
	_, err := s.repo.GetByLogin(login)
	if err != nil && !errors.Is(err, ErrUserByLoginNotExists) {
		return err
	} else if err == nil {
		return ErrUserByLoginAlreadyExists
	}

	encodedPassword, err := s.pwdEncoder.Encode(password)
	if err != nil {
		return err
	}
	user := &User{
		ID:              s.repo.NextID(),
		Login:           login,
		EncodedPassword: encodedPassword,
	}

	return s.repo.Store(user)
}

func NewUserService(repo UserRepository, pwdEncoder PasswordEncoder) *UserService {
	return &UserService{repo, pwdEncoder}
}
