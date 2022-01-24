package app

import (
	"github.com/google/uuid"
)

type UserService interface {
	Create(id uuid.UUID, email, firstName, lastName string) error
	Get(id uuid.UUID) (*User, error)
	Update(id uuid.UUID, email, firstName, lastName string) error
	Delete(id uuid.UUID) error
}

type userService struct {
	repo UserRepository
}

func (s *userService) Create(id uuid.UUID, email, firstName, lastName string) error {
	u := User{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
	return s.repo.Create(&u)
}

func (s *userService) Get(id uuid.UUID) (*User, error) {
	return s.repo.Get(id)
}

func (s *userService) Update(id uuid.UUID, email, firstName, lastName string) error {
	u := User{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
	return s.repo.Update(&u)
}

func (s *userService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo}
}
