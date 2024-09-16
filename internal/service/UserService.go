package service

import (
	"fmt"
	"tender-service/internal/data"
	"tender-service/internal/database/repositories"
)

type UserService struct {
	storage repositories.UserRepository
}

func NewUserService(storage repositories.UserRepository) *UserService {
	return &UserService{storage: storage}
}

func (s *UserService) UserExistByUsername(id string) (bool, error) {
	result, err := s.storage.UserExistByUsername(id)
	if err != nil {
		return false, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return result, nil
}

func (s *UserService) GetUserTendersByUsername(username string) ([]data.Tender, error) {

	userExist, err := s.UserExistByUsername(username)
	if err != nil {
		return nil, err
	}

	if !userExist {
		return nil, UserNotFound
	}

	result, err := s.storage.GetUserTenders(username)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return result, nil
}
