package repositories

import (
	"tender-service/internal/data"
	"tender-service/internal/requests"
)

type BidRepository interface {
	CreateBid(request *requests.BidCreateRequest) (*data.Bid, error)

	GetUserBids(username string) ([]data.Bid, error)

	GetTenderBids(id string) ([]data.Bid, error)

	BidExistById(id string) (bool, error)

	GetBidById(id string) (*data.Bid, error)
	ChangeBidStatus(id, status string) (*data.Bid, error)
}
