package storage

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/userprofile/app"
)

type userProfileRepository struct {
	profiles map[uuid.UUID]*app.UserProfile
}

func (r *userProfileRepository) Get(userID uuid.UUID) (*app.UserProfile, error) {
	profile, ok := r.profiles[userID]
	if !ok {
		return nil, app.ErrUserProfileNotFound
	}
	return profile, nil
}

func (r *userProfileRepository) Store(profile *app.UserProfile) error {
	r.profiles[profile.UserID] = profile
	return nil
}

func NewUserProfileRepository() app.UserProfileRepository {
	return &userProfileRepository{make(map[uuid.UUID]*app.UserProfile)}
}
