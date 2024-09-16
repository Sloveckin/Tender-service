package repositories

import (
	"tender-service/internal/data"
)

type UserRepository interface {
	GetUserTenders(username string) ([]data.Tender, error)

	UserExistByUsername(name string) (bool, error)
}
