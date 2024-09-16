package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"tender-service/internal/data"
	"tender-service/internal/database"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (s *UserRepositoryImpl) GetUserTenders(username string) ([]data.Tender, error) {
	username = PrepareString(username)
	msgError := "database.storage.GetUserTenders"
	op := fmt.Sprintf("%s where creator_username='%s' and last_version='t'", database.SelectTender, username)
	rows, err := s.db.Query(op)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msgError, err)
	}

	var tenders []data.Tender
	for rows.Next() {
		tender, err := GetTenderFromRows(rows, msgError)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}
	return tenders, nil
}

func (s *UserRepositoryImpl) UserExistByUsername(name string) (bool, error) {
	name = PrepareString(name)
	op := fmt.Sprintf("SELECT username FROM Employee WHERE username='%s';", name)
	return s.existBy(op, "database.storage.UserExistByUsername")
}

func (s *UserRepositoryImpl) existBy(op, msgError string) (bool, error) {
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
