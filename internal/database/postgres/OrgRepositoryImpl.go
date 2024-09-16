package postgres

import (
	"database/sql"
	"errors"
	"fmt"
)

type OrgRepositoryImpl struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrgRepositoryImpl {
	return &OrgRepositoryImpl{db}
}

func (s *OrgRepositoryImpl) OrganizationExistsById(organisationId string) (bool, error) {
	op := fmt.Sprintf("SELECT name FROM organization WHERE id='%s';", organisationId)
	return s.existBy(op, "database.storage.OrganizationExistsById")
}

func (s *OrgRepositoryImpl) existBy(op, msgError string) (bool, error) {
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
