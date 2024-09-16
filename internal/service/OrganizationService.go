package service

import (
	"fmt"
	"tender-service/internal/database/repositories"
)

type OrganizationService struct {
	storage repositories.OrgRepository
}

func InitOrganizationService(storage repositories.OrgRepository) *OrganizationService {
	return &OrganizationService{storage: storage}
}

func (s *OrganizationService) ExistById(id string) (bool, error) {
	result, err := s.storage.OrganizationExistsById(id)
	if err != nil {
		return false, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return result, err
}
