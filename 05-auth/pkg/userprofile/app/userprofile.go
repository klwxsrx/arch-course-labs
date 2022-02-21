package app

import (
	"errors"
	"github.com/google/uuid"
)

type UserProfile struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
}

var ErrUserProfileNotFound = errors.New("user profile not found")

type UserProfileRepository interface {
	Get(userID uuid.UUID) (*UserProfile, error)
	Store(profile *UserProfile) error
}

var ErrPermissionDenied = errors.New("can't update another user profile")

type UserProfileService struct {
	repo UserProfileRepository
}

func (s *UserProfileService) Get(userID uuid.UUID, subjectID uuid.UUID) (*UserProfile, error) {
	if userID != subjectID {
		return nil, ErrPermissionDenied
	}

	profile, err := s.repo.Get(userID)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *UserProfileService) Update(userID uuid.UUID, firstName, lastName string, subjectID uuid.UUID) error {
	if userID != subjectID {
		return ErrPermissionDenied
	}

	return s.repo.Store(&UserProfile{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
	})
}

func NewUserProfileService(repo UserProfileRepository) *UserProfileService {
	return &UserProfileService{repo}
}
