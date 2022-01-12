package repository

import "github.com/williepotgieter/nft-maker/pkg/domain/core/adding"

type Database interface {
	adding.Database
}

type DatabasePort struct {
	Add adding.DatabaseService
}

func NewDatabasePort(db Database) *DatabasePort {
	return &DatabasePort{
		Add: adding.NewDBService(db),
	}
}
