package adding

import "github.com/williepotgieter/nft-maker/pkg/domain/models"

type Database interface {
	CreateUser(user models.User) error
}

type DatabaseService interface {
	User(user models.User) error
}

type dbservice struct {
	db Database
}

func NewDBService(db Database) DatabaseService {
	return &dbservice{db}
}
