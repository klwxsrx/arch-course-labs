package storage

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/auth/app"
)

type userRepo struct {
	users map[string]*app.User
}

func (r *userRepo) NextID() uuid.UUID {
	return uuid.New()
}

func (r *userRepo) GetByLogin(login string) (*app.User, error) {
	user, ok := r.users[login]
	if !ok {
		return nil, app.ErrUserByLoginNotExists
	}
	return user, nil
}

func (r *userRepo) Store(user *app.User) error {
	r.users[user.Login] = user
	return nil
}

func NewUserRepository() app.UserRepository {
	return &userRepo{make(map[string]*app.User)}
}
