package service

import (
	"fmt"
	"tender-service/internal/data"
	"tender-service/internal/database/repositories"
	"tender-service/internal/requests"
)

type TenderService struct {
	storage             repositories.TenderRepository
	userService         *UserService
	organizationService *OrganizationService
}

func NewTenderService(storage repositories.TenderRepository, userService *UserService, organizationService *OrganizationService) *TenderService {
	return &TenderService{storage: storage, userService: userService, organizationService: organizationService}
}

func (s *TenderService) ExistById(id string) (bool, error) {
	result, err := s.storage.TenderExistsById(id)
	if err != nil {
		return false, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return result, err
}

func (s *TenderService) ExistByName(name string) (bool, error) {
	result, err := s.storage.TenderExistsByName(name)
	if err != nil {
		return false, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return result, err
}

func (s *TenderService) CreateTender(request *requests.CreateTenderRequest) (*data.Tender, error) {
	userExist, err := s.userService.UserExistByUsername(request.CreatorUserName)
	if err != nil {
		return nil, err
	}
	if !userExist {
		return nil, UserNotFound
	}

	orgExist, err := s.organizationService.ExistById(request.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !orgExist {
		return nil, OrganizationNotFound
	}

	tenderExist, err := s.ExistByName(request.Name)
	if err != nil {
		return nil, err
	}

	if tenderExist {
		return nil, TenderAlreadyExists
	}

	result, err := s.storage.CreateTender(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err.Error())
	}
	return result, err
}

func (s *TenderService) EditTender(id string, request *requests.EditTenderRequest) (*data.Tender, error) {

	tenderExist, err := s.ExistById(id)
	if err != nil {
		return nil, err
	}
	if !tenderExist {
		return nil, TenderNotFound
	}

	tender, err := s.storage.EditTender(id, request)
	return tender, err
}

func (s *TenderService) GetTenderById(id string) (*data.Tender, error) {
	exist, err := s.ExistById(id)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, TenderNotFound
	}

	return s.storage.GetTenderById(id)
}

func (s *TenderService) GetTenders() ([]data.Tender, error) {
	tenders, err := s.storage.GetTenders()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err.Error())
	}
	return tenders, err
}

func (s *TenderService) RollbackTender(id string, version int64) (*data.Tender, error) {
	tender, err := s.GetTenderById(id)
	if err != nil {
		return nil, err
	}

	if version <= 0 {
		return nil, fmt.Errorf("%w: Version not positive number", NotCorrectParams)
	}

	if version > tender.Version {
		return nil, fmt.Errorf("%w: Tender version is greater than the latest version", NotCorrectParams)
	}

	return s.storage.RollbackTender(id, version)
}
