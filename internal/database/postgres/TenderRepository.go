package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"tender-service/internal/data"
	"tender-service/internal/database"
	requests2 "tender-service/internal/requests"
)

type TenderRepositoryImpl struct {
	db *sql.DB
}

func NewTenderRepository(db *sql.DB) *TenderRepositoryImpl {
	return &TenderRepositoryImpl{db: db}
}

func (s *TenderRepositoryImpl) CreateTender(request *requests2.CreateTenderRequest) (*data.Tender, error) {
	organizationId := request.OrganizationId
	op := fmt.Sprintf("SELECT id from Employee where username='%s';", PrepareString(request.CreatorUserName))

	var userId string
	err := s.db.QueryRow(op).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.CreateTender", err)
	}

	err = s.CreateOrganizationResponsible(userId, organizationId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.CreateTender", err)
	}

	opTender := fmt.Sprintf("%s values (%d, '%s', '%s', '%s', '%s', '%s', '%s') RETURNING %s;",
		database.InsertTender,
		1,
		PrepareString(request.Name),
		PrepareString(request.Description),
		request.Status,
		request.ServiceType,
		request.OrganizationId,
		PrepareString(request.CreatorUserName),
		database.CollectTenderFromRow)

	var tender data.Tender
	row := s.db.QueryRow(opTender)
	err = getTenderFromRow(row, &tender)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.CreateTender", err)
	}
	return &tender, nil
}

func (s *TenderRepositoryImpl) EditTender(id string, request *requests2.EditTenderRequest) (*data.Tender, error) {

	tender, err := s.GetTenderById(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.EditTender", err)
	}

	op := fmt.Sprintf("UPDATE Tenders set last_version='f' where id='%s';", id)
	_, err = s.db.Exec(op)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.EditTender", err)
	}

	var newName string
	if request.Name != "" {
		newName = request.Name
	} else {
		newName = tender.Name
	}

	var newDescription string
	if request.Description != "" {
		newDescription = request.Description
	} else {
		newDescription = tender.Description
	}

	var newStatus string
	if request.Status != "" {
		newStatus = request.Status
	} else {
		newStatus = tender.Status
	}

	op = fmt.Sprintf("%s values ('%s', %d, '%s', '%s', '%s', '%s', '%s', '%s') RETURNING %s;",
		database.InsertTenderVersion, tender.Id, tender.Version+1, PrepareString(newName), PrepareString(newDescription), newStatus, tender.ServiceType, tender.OrganisationId, PrepareString(tender.CreatorUserName), database.CollectTenderFromRow)

	var newTender data.Tender
	row := s.db.QueryRow(op)
	err = getTenderFromRow(row, &newTender)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.CreateTender", err)
	}

	return &newTender, nil
}

func (s *TenderRepositoryImpl) getTenderBy(op, msgError string) (*data.Tender, error) {
	var tender data.Tender
	row := s.db.QueryRow(op)
	err := getTenderFromRow(row, &tender)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", msgError, err)
	}
	return &tender, nil
}

func (s *TenderRepositoryImpl) CreateOrganizationResponsible(userId, organizationId string) error {
	op := fmt.Sprintf("INSERT INTO Organization_responsible (user_id, organization_id) VALUES ('%s', '%s');", userId, organizationId)
	_, err := s.db.Exec(op)
	if err != nil {
		return fmt.Errorf("%s: %w", "database.storage.CreateOrganizationResponsible", err)
	}
	return nil
}

func (s *TenderRepositoryImpl) RollbackTender(id string, version int64) (*data.Tender, error) {
	dropOp := fmt.Sprintf("DELETE from Tenders where id='%s' and version > %d;", id, version)
	_, err := s.db.Exec(dropOp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.RollbackTender-Drop", err)
	}

	updateOp := fmt.Sprintf("UPDATE Tenders set last_version='t' where id='%s' and version='%d' RETURNING %s;", id, version, database.CollectTenderFromRow)
	var tender data.Tender
	row := s.db.QueryRow(updateOp)
	err = getTenderFromRow(row, &tender)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "database.storage.RollbackTender-Update", err)
	}

	return &tender, nil
}

func (s *TenderRepositoryImpl) GetTenderById(id string) (*data.Tender, error) {
	op := fmt.Sprintf("%s where id='%s' and last_version='t'", database.SelectTender, id)
	return s.getTenderBy(op, "database.storage.GetTenderById")
}

func (s *TenderRepositoryImpl) TenderExistsByName(name string) (bool, error) {
	name = PrepareString(name)
	op := fmt.Sprintf("SELECT name FROM Tenders where name='%s';", name)
	return s.existBy(op, "database.storage.TenderExistsByName")
}

func (s *TenderRepositoryImpl) TenderExistsById(id string) (bool, error) {
	op := fmt.Sprintf("SELECT name FROM Tenders where id='%s';", id)
	return s.existBy(op, "database.storage.TenderExistsById")
}

func (s *TenderRepositoryImpl) existBy(op, msgError string) (bool, error) {
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

func (s *TenderRepositoryImpl) GetTenders() ([]data.Tender, error) {
	msgError := "database.storage.GetTenders"
	op := fmt.Sprintf("%s where last_version='t';", database.SelectTender)
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

func GetTenderFromRows(rows *sql.Rows, msgError string) (data.Tender, error) {
	var tender data.Tender
	err := rows.Scan(&tender.Id, &tender.Version, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.OrganisationId, &tender.CreatorUserName, &tender.CreatedAt)
	if err != nil {
		return tender, fmt.Errorf("%s: %w", msgError, err)
	}
	return tender, nil
}

func getTenderFromRow(rows *sql.Row, tender *data.Tender) error {
	return rows.Scan(
		&tender.Id,
		&tender.Version,
		&tender.Name,
		&tender.Description,
		&tender.Status,
		&tender.ServiceType,
		&tender.OrganisationId,
		&tender.CreatorUserName,
		&tender.CreatedAt)
}
