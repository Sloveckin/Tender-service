package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"tender-service/internal/data"
	"tender-service/internal/database"
	"tender-service/internal/requests"
	"tender-service/internal/service"
)

type BidRepository struct {
	db *sql.DB
}

func NewBidRepository(db *sql.DB) *BidRepository {
	return &BidRepository{db}
}

func (s *BidRepository) CreateBid(request *requests.BidCreateRequest) (*data.Bid, error) {
	op := fmt.Sprintf("%s values ('%s', '%s', '%s', '%s', '%s') RETURNING %s;",
		database.InsertBid,
		request.TenderId,
		PrepareString(request.Name),
		PrepareString(request.Description),
		PrepareString(request.CreatorUserName),
		request.OrganizationId,
		database.CollectBidFromRow)

	var bid data.Bid
	row := s.db.QueryRow(op)
	err := row.Scan(&bid.Id, &bid.TenderId, &bid.Name, &bid.Description, &bid.Status, &bid.CreatorUserName, &bid.OrganizationId, &bid.CreatedTime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.CreateBid", err)
	}
	return &bid, nil
}

func (s *BidRepository) BidExistById(id string) (bool, error) {
	op := fmt.Sprintf("SELECT name FROM Bids WHERE id = '%s';", id)
	msgError := "database.TenderRepository.BidExistById"
	return s.existBy(op, msgError)
}

func (s *BidRepository) existBy(op, msgError string) (bool, error) {
	var t string
	err := s.db.QueryRow(op).Scan(&t)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", msgError, err)
	}
	return true, nil
}

func (s *BidRepository) GetUserBids(username string) ([]data.Bid, error) {
	username = PrepareString(username)
	msgError := "database.storage.GetUserBids"
	where := fmt.Sprintf("creator_username='%s'", username)
	return s.getBidsBy(where, msgError)
}

func (s *BidRepository) GetTenderBids(id string) ([]data.Bid, error) {
	msgError := "database.storage.TenderBids"
	where := fmt.Sprintf("tender_id='%s'", id)
	return s.getBidsBy(where, msgError)
}

func (s *BidRepository) GetBidById(id string) (*data.Bid, error) {
	op := fmt.Sprintf("%s where id='%s';", database.SelectBid, id)
	row := s.db.QueryRow(op)
	var bid data.Bid
	err := row.Scan(&bid.Id, &bid.TenderId, &bid.Name, &bid.Description, &bid.Status, &bid.CreatorUserName, &bid.OrganizationId)
	if err != nil {
		return &bid, fmt.Errorf("%s: %w", service.NotExpectedErrorFromStorage, err)
	}
	return &bid, nil
}

func (s *BidRepository) getBidsBy(where, msgError string) ([]data.Bid, error) {
	op := fmt.Sprintf("%s where %s;", database.SelectBid, where)
	rows, err := s.db.Query(op)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msgError, err)
	}

	var bids []data.Bid
	for rows.Next() {
		bid, err := getBidFromRows(rows, msgError)
		if err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

func (s *BidRepository) ChangeBidStatus(id, status string) (*data.Bid, error) {
	panic("implement me")
}

func getBidFromRows(rows *sql.Rows, msgError string) (data.Bid, error) {
	var bid data.Bid
	err := rows.Scan(&bid.Id, &bid.TenderId, &bid.Name, &bid.Description, &bid.Status, &bid.CreatorUserName, &bid.OrganizationId)
	if err != nil {
		return bid, fmt.Errorf("%s: %w", msgError, err)
	}
	return bid, nil
}
