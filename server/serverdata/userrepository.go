package serverdata

import (
	"errors"
	"tchat/internal/types"
)

type UserRepository struct {
	users map[string]types.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]types.User),
	}
}

func (r *UserRepository) AddUser(user types.User) error {
	if _, ok := r.users[user.UserID]; ok {
		return errors.New("user already exists")
	}

	r.users[user.UserID] = user
	return nil
}

func (r *UserRepository) GetUser(userID string) *types.User {
	user, ok := r.users[userID]
	if !ok {
		return nil
	}

	return &user
}
