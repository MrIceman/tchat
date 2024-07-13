package serverdomain

import (
	"tchat/internal/types"
	"tchat/server/serverdata"
	"time"
)

type UserService struct {
	repository *serverdata.UserRepository
}

func NewService(repository *serverdata.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) SignInUser(userID string) error {
	return s.repository.AddUser(types.User{
		UserID:     userID,
		LoggedInAt: time.Now(),
	})
}

func (s *UserService) DoesUserExist(userID string) bool {
	return s.repository.GetUser(userID) != nil
}
