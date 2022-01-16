package reading

import "github.com/williepotgieter/nft-maker/pkg/domain/models"

// Database service
type Database interface {
	GetAllUsers() ([]models.User, error)
}

type DatabaseService interface {
	AllUsers() ([]models.User, error)
}

type dbservice struct {
	db Database
}

func NewDBService(db Database) DatabaseService {
	return &dbservice{db}
}
