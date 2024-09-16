package repositories

import (
	"tender-service/internal/data"
	"tender-service/internal/requests"
)

type TenderRepository interface {
	CreateTender(request *requests.CreateTenderRequest) (*data.Tender, error)

	EditTender(id string, request *requests.EditTenderRequest) (*data.Tender, error)

	CreateOrganizationResponsible(userId, organizationId string) error

	RollbackTender(id string, version int64) (*data.Tender, error)

	GetTenderById(id string) (*data.Tender, error)

	TenderExistsByName(name string) (bool, error)

	TenderExistsById(id string) (bool, error)

	GetTenders() ([]data.Tender, error)
}
