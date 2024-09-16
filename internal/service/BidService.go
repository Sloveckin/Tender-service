package service

import (
	"fmt"
	"tender-service/internal/data"
	"tender-service/internal/database/repositories"
	"tender-service/internal/requests"
)

type BidService struct {
	storage             repositories.BidRepository
	userService         *UserService
	organizationService *OrganizationService
	tenderService       *TenderService
}

func InitBidService(storage repositories.BidRepository, userService *UserService, orgService *OrganizationService, tenderService *TenderService) *BidService {
	return &BidService{storage: storage, userService: userService, organizationService: orgService, tenderService: tenderService}
}

func (s *BidService) CreateBid(request *requests.BidCreateRequest) (*data.Bid, error) {

	tenderExist, err := s.tenderService.ExistById(request.TenderId)
	if err != nil {
		return nil, err
	}

	if !tenderExist {
		return nil, TenderNotFound
	}

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

	bid, err := s.storage.CreateBid(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return bid, nil
}

func (s *BidService) GetTenderBids(id string) ([]data.Bid, error) {
	exist, err := s.tenderService.ExistById(id)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, TenderNotFound
	}

	bids, err := s.storage.GetTenderBids(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return bids, nil
}

func (s *BidService) GetUserBids(username string) ([]data.Bid, error) {
	exist, err := s.userService.UserExistByUsername(username)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, UserNotFound
	}

	bids, err := s.storage.GetUserBids(username)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return bids, err
}

func (s *BidService) BidExistById(id string) (bool, error) {
	exist, err := s.storage.BidExistById(id)
	if err != nil {
		return false, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return exist, nil
}

func (s *BidService) GetBidById(id string) (*data.Bid, error) {
	exist, err := s.BidExistById(id)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, BidNotFound
	}

	bid, err := s.storage.GetBidById(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return bid, nil
}

func (s *BidService) ChangeBidStatusById(id, status string) (*data.Bid, error) {
	exist, err := s.BidExistById(id)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, BidNotFound
	}

	bid, err := s.storage.ChangeBidStatus(id, status)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", NotExpectedErrorFromStorage, err)
	}
	return bid, nil
}
